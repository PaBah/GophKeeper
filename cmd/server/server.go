package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/PaBah/GophKeeper/internal/auth"
	"github.com/PaBah/GophKeeper/internal/config"
	pb "github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
	"github.com/PaBah/GophKeeper/internal/models"
	"github.com/PaBah/GophKeeper/internal/storage"
	"github.com/PaBah/GophKeeper/internal/utils"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb.UnimplementedGophKeeperServiceServer
	config      *config.ServerConfig
	storage     storage.Repository
	minioClient *minio.Client

	syncClients map[string]map[string]pb.GophKeeperService_SubscribeToChangesServer
	rwMutex     *sync.RWMutex
}

// SignIn - handler for Sign In
func (s *GrpcServer) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.SignInResponse, error) {
	response := &pb.SignInResponse{}
	user, err := s.storage.AuthorizeUser(ctx, in.Email)

	if err != nil || !utils.CheckPasswordHash(user.Password, in.Password) {
		return response, status.Errorf(codes.Unavailable, "User with such credentials can not be logined")
	}

	sessionID := uuid.New().String()
	JWTToken, err := auth.BuildJWTString(user.ID, sessionID, s.config.Secret)
	if err != nil {
		return response, status.Errorf(codes.Internal, "Can not build auth token")
	}

	response.Token = JWTToken
	return response, nil
}

// SignUp - handler for Sign Up
func (s *GrpcServer) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	response := &pb.SignUpResponse{}

	user := models.NewUser(in.Email, in.Password)
	createdUser, err := s.storage.CreateUser(ctx, user)

	if errors.Is(err, storage.ErrAlreadyExists) {
		return response, status.Errorf(codes.InvalidArgument, "User with such email already exists")
	}

	sessionID := uuid.New().String()
	JWTToken, err := auth.BuildJWTString(createdUser.ID, sessionID, s.config.Secret)
	if err != nil {
		return response, status.Errorf(codes.Internal, "JWT token can not be built")
	}

	err = s.minioClient.MakeBucket(context.Background(), createdUser.ID, minio.MakeBucketOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	response.Token = JWTToken
	return response, nil
}

// CreateCredentials - handler for creating Credentials records in DB
func (s *GrpcServer) CreateCredentials(ctx context.Context, in *pb.CreateCredentialsRequest) (*pb.CreateCredentialsResponse, error) {
	response := &pb.CreateCredentialsResponse{}

	credentials := models.NewCredentials(in.ServiceName, in.Identity, in.Password)
	createdCredentials, err := s.storage.CreateCredentials(ctx, credentials)

	if errors.Is(err, storage.ErrAlreadyExists) {
		return response, status.Errorf(codes.InvalidArgument, "User already created credentials with such service name and identity")
	}

	s.SendNotifications(ctx, 0, createdCredentials.ID)
	response.Id = createdCredentials.ID
	response.ServiceName = createdCredentials.ServiceName
	response.UploadedAt = createdCredentials.UploadedAt.Format(time.RFC3339)

	return response, nil
}

// GetCredentials - handler for get credentials stored by user
func (s *GrpcServer) GetCredentials(ctx context.Context, in *pb.GetCredentialsRequest) (*pb.GetCredentialsResponse, error) {
	response := &pb.GetCredentialsResponse{}

	credentials, err := s.storage.GetCredentials(ctx)

	if err != nil {
		return response, status.Errorf(codes.InvalidArgument, "credentials can not be retrieved")
	}

	for _, credentialSet := range credentials {
		response.Credentials = append(response.Credentials, &pb.GetCredentialsResponse_Credential{
			Id:          credentialSet.ID,
			ServiceName: credentialSet.ServiceName,
			Identity:    credentialSet.Identity,
			Password:    credentialSet.Password,
			UploadedAt:  credentialSet.UploadedAt.Format(time.RFC3339),
		})
	}
	return response, nil
}

// UpdateCredentials - handler for get credentials stored by user
func (s *GrpcServer) UpdateCredentials(ctx context.Context, in *pb.UpdateCredentialsRequest) (*pb.UpdateCredentialsResponse, error) {
	response := &pb.UpdateCredentialsResponse{}

	createdCredentials, err := s.storage.UpdateCredentials(ctx, models.Credentials{
		ID:          in.Id,
		ServiceName: in.ServiceName,
		Identity:    in.Identity,
		Password:    in.Password,
	})

	if err != nil {
		return response, status.Errorf(codes.InvalidArgument, "credentials can not be updated")
	}
	s.SendNotifications(ctx, 0, createdCredentials.ID)
	response.Id = createdCredentials.ID
	response.ServiceName = createdCredentials.ServiceName
	response.UploadedAt = createdCredentials.UploadedAt.Format(time.RFC3339)

	return response, nil
}

// DeleteCredentials - handler for get credentials stored by user
func (s *GrpcServer) DeleteCredentials(ctx context.Context, in *pb.DeleteCredentialsRequest) (*pb.DeleteCredentialsResponse, error) {
	response := &pb.DeleteCredentialsResponse{}

	err := s.storage.DeleteCredentials(ctx, in.Id)

	if err != nil {
		return response, status.Errorf(codes.InvalidArgument, "credentials can not be deleted")
	}
	s.SendNotifications(ctx, 0, in.Id)
	return response, nil
}

// CreateCard - handler for creating Card records in DB
func (s *GrpcServer) CreateCard(ctx context.Context, in *pb.CreateCardRequest) (*pb.CreateCardResponse, error) {
	response := &pb.CreateCardResponse{}

	if utils.ValidateLuhn(in.Number) != nil {
		return response, status.Errorf(codes.InvalidArgument, "invalid card number")
	}

	card := models.NewCard(in.Number, in.ExpirationDate, in.HolderName, in.Cvv)
	createdCard, err := s.storage.CreateCard(ctx, card)

	if err != nil {
		return response, status.Errorf(codes.InvalidArgument, "card can not be created")
	}
	s.SendNotifications(ctx, 1, createdCard.ID)
	response.LastDigits = createdCard.Number[12:]
	response.ExpirationDate = createdCard.ExpirationDate
	response.UploadedAt = createdCard.UploadedAt.Format(time.RFC3339)

	return response, nil
}

// GetCards - handler for get Cards stored by user
func (s *GrpcServer) GetCards(ctx context.Context, in *pb.GetCardsRequest) (*pb.GetCardsResponse, error) {
	response := &pb.GetCardsResponse{}

	cards, err := s.storage.GetCards(ctx)
	if err != nil {
		return response, status.Errorf(codes.InvalidArgument, "cards can not be retrieved")
	}

	for _, card := range cards {
		response.Cards = append(response.Cards, &pb.GetCardsResponse_Card{
			Id:             card.ID,
			Number:         card.Number,
			ExpirationDate: card.ExpirationDate,
			HolderName:     card.HolderName,
			Cvv:            card.CVV,
			UploadedAt:     card.UploadedAt.Format(time.RFC3339),
		})
	}
	return response, nil
}

// UpdateCard - handler for update Card stored by user
func (s *GrpcServer) UpdateCard(ctx context.Context, in *pb.UpdateCardRequest) (*pb.UpdateCardResponse, error) {
	response := &pb.UpdateCardResponse{}

	if utils.ValidateLuhn(in.Number) != nil {
		return response, status.Errorf(codes.InvalidArgument, "invalid card number")
	}

	card, err := s.storage.UpdateCard(ctx, models.Card{
		ID:             in.Id,
		Number:         in.Number,
		ExpirationDate: in.ExpirationDate,
		HolderName:     in.HolderName,
		CVV:            in.Cvv,
	})

	if err != nil {
		return response, status.Errorf(codes.InvalidArgument, "card can not be updated")
	}
	s.SendNotifications(ctx, 1, card.ID)
	response.LastDigits = card.Number[12:]
	response.ExpirationDate = card.ExpirationDate
	response.UploadedAt = card.UploadedAt.Format(time.RFC3339)

	return response, nil
}

// DeleteCard - handler for deletion of user's Card
func (s *GrpcServer) DeleteCard(ctx context.Context, in *pb.DeleteCardRequest) (*pb.DeleteCardResponse, error) {
	response := &pb.DeleteCardResponse{}

	err := s.storage.DeleteCard(ctx, in.Id)

	if err != nil {
		return response, status.Errorf(codes.InvalidArgument, "card can not be deleted")
	}
	s.SendNotifications(ctx, 1, in.Id)
	return response, nil
}

// SubscribeToChanges - stream changes to clients
func (s *GrpcServer) SubscribeToChanges(in *pb.SubscribeToChangesRequest, stream pb.GophKeeperService_SubscribeToChangesServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	userID, _ := ctx.Value(config.USERIDCONTEXTKEY).(string)
	sessionID, _ := ctx.Value(config.SESSIONIDCONTEXTKEY).(string)
	defer cancel()
	s.rwMutex.Lock()
	if len(s.syncClients[userID]) == 0 {
		s.syncClients[userID] = make(map[string]pb.GophKeeperService_SubscribeToChangesServer)
	}
	s.syncClients[userID][sessionID] = stream
	s.rwMutex.Unlock()
	for {
		time.Sleep(time.Minute)
	}
}

// SendNotifications - stream all user session with update
func (s *GrpcServer) SendNotifications(ctx context.Context, resource int, ID string) {
	sessionID, _ := ctx.Value(config.SESSIONIDCONTEXTKEY).(string)
	userID, _ := ctx.Value(config.USERIDCONTEXTKEY).(string)
	s.rwMutex.Lock()
	for session, client := range s.syncClients[userID] {
		if session != sessionID {
			_ = client.Send(&pb.SubscribeToChangesResponse{
				Source: int32(resource),
				Id:     ID,
			})
		}
	}
	s.rwMutex.Unlock()
}

func (s *GrpcServer) UploadFile(stream pb.GophKeeperService_UploadFileServer) (err error) {
	var objectName string
	var fileData []byte
	userID, _ := stream.Context().Value(config.USERIDCONTEXTKEY).(string)

	for {
		var in *pb.UploadFileRequest
		in, err = stream.Recv()
		if err == io.EOF {
			reader := io.NopCloser(bytes.NewReader(fileData))
			_, err = s.minioClient.PutObject(context.Background(), userID, objectName, reader, int64(len(fileData)), minio.PutObjectOptions{})
			if err != nil {
				return
			}
			s.SendNotifications(stream.Context(), 2, objectName)
			return stream.Send(&pb.UploadFileResponse{
				Message: "File uploaded successfully",
				Success: true,
			})
		}
		if err != nil {
			return
		}

		fileData = append(fileData, in.Data...)
		objectName = in.Filename
	}
}

func (s *GrpcServer) GetFiles(ctx context.Context, in *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	response := &pb.GetFilesResponse{}

	objectCh := s.minioClient.ListObjects(ctx, ctx.Value(config.USERIDCONTEXTKEY).(string), minio.ListObjectsOptions{})
	for object := range objectCh {
		response.Files = append(response.Files, &pb.GetFilesResponse_File{
			Name:       object.Key,
			Size:       utils.HumanReadableSize(uint64(object.Size)),
			UploadedAt: object.LastModified.Format(time.RFC3339),
		})
	}
	return response, nil
}

func (s *GrpcServer) DeleteFile(ctx context.Context, in *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	response := &pb.DeleteFileResponse{}

	err := s.minioClient.RemoveObject(ctx, ctx.Value(config.USERIDCONTEXTKEY).(string), in.Name, minio.RemoveObjectOptions{})
	s.SendNotifications(ctx, 2, in.Name)
	return response, err
}

func (s *GrpcServer) DownloadFile(in *pb.DownloadFileRequest, stream pb.GophKeeperService_DownloadFileServer) error {
	object, err := s.minioClient.GetObject(
		context.Background(), stream.Context().Value(config.USERIDCONTEXTKEY).(string),
		in.Name, minio.GetObjectOptions{},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to open object: %v", err)
	}
	defer object.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := object.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "block can not be read: %v", err)
		}

		if err := stream.Send(&pb.DownloadFileResponse{Data: buffer[:n]}); err != nil {
			return status.Errorf(codes.Internal, "failed to send chunk: %v", err)
		}
	}

	return nil
}

// NewGrpcServer - creates new gRPC server instance
func NewGrpcServer(config *config.ServerConfig, storage storage.Repository) *GrpcServer {
	// Настраиваем соединение с MinIO
	minioClient, err := minio.New("127.0.0.1:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("admin", "password123", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	s := GrpcServer{
		config:      config,
		storage:     storage,
		minioClient: minioClient,
		syncClients: make(map[string]map[string]pb.GophKeeperService_SubscribeToChangesServer),
		rwMutex:     &sync.RWMutex{},
	}
	return &s
}
