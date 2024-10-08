package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	pb "github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
	"github.com/PaBah/GophKeeper/internal/logger"
	"github.com/PaBah/GophKeeper/internal/models"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type Client interface {
	SignUp(email, password string) error
	SignIn(email, password string) error
}

type ClientService struct {
	client        pb.GophKeeperServiceClient
	token         string
	sessionID     string
	serverAddress string
	conn          *grpc.ClientConn
	isAvailable   bool
}

type GRPCClientProvider interface {
	SignUp(email, password string) error
	SignIn(email, password string) error
	CreateCredentials(ctx context.Context, serviceName, identity, password string) error
	GetCredentials(ctx context.Context) (credentials []models.Credentials, err error)
	UpdateCredentials(ctx context.Context, credentials models.Credentials) (updatedCredentials models.Credentials, err error)
	DeleteCredentials(ctx context.Context, credentialsID string) (err error)
	CreateCard(ctx context.Context, number, expirationDate, holderName, cvv string) error
	GetCards(ctx context.Context) (cards []models.Card, err error)
	GetFiles(ctx context.Context) (files []models.File, err error)
	DeleteFile(ctx context.Context, name string) (err error)
	UpdateCards(ctx context.Context, card models.Card) (updatedCard models.Card, err error)
	DeleteCard(ctx context.Context, cardID string) (err error)
	UploadFile(ctx context.Context, filePath string)
	DownloadsFile(ctx context.Context, name string)
	SubscribeToChanges(ctx context.Context) (grpc.ServerStreamingClient[pb.SubscribeToChangesResponse], error)
	TryToConnect() bool
}

func NewClientService(serverAddress string) ClientService {
	return ClientService{
		serverAddress: serverAddress,
	}
}

// SignUp registers a new user with the provided email and password, and stores the authentication token.
func (c *ClientService) SignUp(email, password string) error {
	resp, err := c.client.SignUp(context.Background(), &pb.SignUpRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("SignUp: %w", err)
	}
	c.token = resp.GetToken()

	return nil
}

func (c *ClientService) SignIn(email, password string) error {
	resp, err := c.client.SignIn(context.Background(), &pb.SignInRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("SignIn: %w", err)
	}
	c.token = resp.GetToken()

	return nil
}

func (c *ClientService) CreateCredentials(ctx context.Context, serviceName, identity, password string) error {
	_, err := c.client.CreateCredentials(c.getCtx(ctx, c.token), &pb.CreateCredentialsRequest{
		ServiceName: serviceName,
		Identity:    identity,
		Password:    password,
	})
	if err != nil {
		return fmt.Errorf("CreateCredentials: %w", err)
	}
	return nil
}

func (c *ClientService) GetCredentials(ctx context.Context) (credentials []models.Credentials, err error) {
	resp, err := c.client.GetCredentials(c.getCtx(ctx, c.token), &pb.GetCredentialsRequest{})
	if err != nil {
		err = fmt.Errorf("CreateCredentials: %w", err)
		return
	}
	for _, cred := range resp.Credentials {
		uploadedAt, _ := time.Parse(time.RFC3339, cred.UploadedAt)
		credentials = append(credentials, models.Credentials{
			ID:          cred.Id,
			ServiceName: cred.ServiceName,
			Identity:    cred.Identity,
			Password:    cred.Password,
			UploadedAt:  uploadedAt,
		})
	}
	return
}

func (c *ClientService) UpdateCredentials(ctx context.Context, credentials models.Credentials) (updatedCredentials models.Credentials, err error) {
	response, err := c.client.UpdateCredentials(c.getCtx(ctx, c.token), &pb.UpdateCredentialsRequest{
		Id:          credentials.ID,
		ServiceName: credentials.ServiceName,
		Identity:    credentials.Identity,
		Password:    credentials.Password,
	})
	if err != nil {
		err = fmt.Errorf("UpdateCredentials: %w", err)
		return
	}
	updatedCredentials = credentials
	updatedCredentials.ServiceName = response.ServiceName
	uploadedAt, _ := time.Parse(time.RFC3339, response.UploadedAt)
	updatedCredentials.UploadedAt = uploadedAt

	return
}

func (c *ClientService) DeleteCredentials(ctx context.Context, credentialsID string) (err error) {
	_, err = c.client.DeleteCredentials(c.getCtx(ctx, c.token), &pb.DeleteCredentialsRequest{
		Id: credentialsID,
	})
	if err != nil {
		err = fmt.Errorf("DeleteCredentials: %w", err)
		return
	}

	return
}

func (c *ClientService) CreateCard(ctx context.Context, number, expirationDate, holderName, cvv string) error {
	_, err := c.client.CreateCard(c.getCtx(ctx, c.token), &pb.CreateCardRequest{
		Number:         number,
		ExpirationDate: expirationDate,
		HolderName:     holderName,
		Cvv:            cvv,
	})

	if err != nil {
		return fmt.Errorf("CreateCard: %w", err)
	}
	return nil
}

func (c *ClientService) GetCards(ctx context.Context) (cards []models.Card, err error) {
	resp, err := c.client.GetCards(c.getCtx(ctx, c.token), &pb.GetCardsRequest{})
	if err != nil {
		err = fmt.Errorf("GetCards: %w", err)
		return
	}
	for _, card := range resp.Cards {
		uploadedAt, _ := time.Parse(time.RFC3339, card.UploadedAt)
		cards = append(cards, models.Card{
			ID:             card.Id,
			Number:         card.Number,
			ExpirationDate: card.ExpirationDate,
			HolderName:     card.HolderName,
			CVV:            card.Cvv,
			UploadedAt:     uploadedAt,
		})
	}
	return
}

func (c *ClientService) GetFiles(ctx context.Context) (files []models.File, err error) {
	resp, err := c.client.GetFiles(c.getCtx(ctx, c.token), &pb.GetFilesRequest{})
	if err != nil {
		err = fmt.Errorf("GetFiles: %w", err)
		return
	}
	for _, card := range resp.Files {
		uploadedAt, _ := time.Parse(time.RFC3339, card.UploadedAt)
		files = append(files, models.File{
			Name:       card.Name,
			Size:       card.Size,
			UploadedAt: uploadedAt,
		})
	}
	return
}

func (c *ClientService) DeleteFile(ctx context.Context, name string) (err error) {
	_, err = c.client.DeleteFile(c.getCtx(ctx, c.token), &pb.DeleteFileRequest{
		Name: name,
	})

	if err != nil {
		err = fmt.Errorf("DeleteFile: %w", err)
	}

	return
}

func (c *ClientService) UpdateCards(ctx context.Context, card models.Card) (updatedCard models.Card, err error) {
	response, err := c.client.UpdateCard(c.getCtx(ctx, c.token), &pb.UpdateCardRequest{
		Id:             card.ID,
		Number:         card.Number,
		ExpirationDate: card.ExpirationDate,
		HolderName:     card.HolderName,
		Cvv:            card.CVV,
	})
	if err != nil {
		err = fmt.Errorf("UpdateCards: %w", err)
		return
	}
	updatedCard = card
	updatedCard.ExpirationDate = response.ExpirationDate
	uploadedAt, _ := time.Parse(time.RFC3339, response.UploadedAt)
	updatedCard.UploadedAt = uploadedAt

	return
}

func (c *ClientService) DeleteCard(ctx context.Context, cardID string) (err error) {
	_, err = c.client.DeleteCard(c.getCtx(ctx, c.token), &pb.DeleteCardRequest{
		Id: cardID,
	})
	if err != nil {
		err = fmt.Errorf("DeleteCard: %w", err)
		return
	}

	return
}

func (c *ClientService) UploadFile(ctx context.Context, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Log().Error("could not open file:", zap.Error(err))
		return
	}
	defer file.Close()

	stream, err := c.client.UploadFile(c.getCtx(ctx, c.token))
	if err != nil {
		logger.Log().Error("could not upload file:", zap.Error(err))
		return
	}

	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Log().Error("could not read file:", zap.Error(err))
			return
		}

		if err := stream.Send(&pb.UploadFileRequest{
			Data:     buffer[:n],
			Filename: filepath.Base(filePath),
		}); err != nil {
			logger.Log().Error("could not send chunk:", zap.Error(err))
			return
		}
	}
	if err := stream.CloseSend(); err != nil {
		logger.Log().Error("could not close stream:", zap.Error(err))
		return
	}

	resp, err := stream.Recv()
	if err != nil {
		logger.Log().Error("could not receive response:", zap.Error(err))
		return
	}
	logger.Log().Info("Response from server:", zap.Any("response", resp))
}

func (c *ClientService) DownloadsFile(ctx context.Context, name string) {
	req := &pb.DownloadFileRequest{Name: name}

	stream, err := c.client.DownloadFile(c.getCtx(ctx, c.token), req)
	if err != nil {
		logger.Log().Error("error downloading file: ", zap.Error(err))
		return
	}

	localFile, err := os.Create(name)
	if err != nil {
		logger.Log().Error("could not create local file: ", zap.Error(err))
		return
	}
	defer localFile.Close()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Log().Error("error receiving chunk: ", zap.Error(err))
			return
		}

		if _, err := localFile.Write(resp.Data); err != nil {
			logger.Log().Error("error writing to local file: ", zap.Error(err))
			return
		}
	}

	logger.Log().Info("File downloaded successfully.")
}

func (c *ClientService) SubscribeToChanges(ctx context.Context) (grpc.ServerStreamingClient[pb.SubscribeToChangesResponse], error) {
	return c.client.SubscribeToChanges(c.getCtx(ctx, c.token), &pb.SubscribeToChangesRequest{})
}

// TryToConnect attempts to establish a connection with the gRPC server.
// It sets up the connection and checks the server's availability.
func (c *ClientService) TryToConnect() bool {
	conn, err := grpc.NewClient(c.serverAddress,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	if err != nil {
		logger.Log().Error("failed connect to server", zap.Error(err))
		return false
	}

	c.conn = conn
	c.client = pb.NewGophKeeperServiceClient(conn)

	c.isAvailable = true

	return true
}

func (c *ClientService) getCtx(ctx context.Context, jwt string) context.Context {
	md := metadata.New(map[string]string{
		"authorization": jwt,
		"session":       c.sessionID,
	})

	newCtx := metadata.NewOutgoingContext(ctx, md)

	return newCtx
}
