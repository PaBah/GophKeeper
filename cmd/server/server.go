package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PaBah/GophKeeper/internal/auth"
	"github.com/PaBah/GophKeeper/internal/config"
	pb "github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
	"github.com/PaBah/GophKeeper/internal/models"
	"github.com/PaBah/GophKeeper/internal/storage"
	"github.com/PaBah/GophKeeper/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	pb.UnimplementedGophKeeperServiceServer
	logger  *zap.Logger
	config  *config.ServerConfig
	storage storage.Repository
}

// SignIn - handler for Sign In
func (s *GrpcServer) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.SignInResponse, error) {
	response := &pb.SignInResponse{}
	user, err := s.storage.AuthorizeUser(ctx, in.Email)

	if err != nil || !utils.CheckPasswordHash(user.Password, in.Password) {
		return response, status.Errorf(codes.Unavailable, "User with such credentials can not be logined")
	}

	JWTToken, err := auth.BuildJWTString(user.ID, s.config.Secret)
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
	fmt.Println(createdUser)
	JWTToken, err := auth.BuildJWTString(createdUser.ID, s.config.Secret)
	if err != nil {
		return response, status.Errorf(codes.Unauthenticated, err.Error())
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

	if err != nil {
		return response, status.Errorf(codes.Unauthenticated, err.Error())
	}

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
		return response, status.Errorf(codes.InvalidArgument, err.Error())
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
		return response, status.Errorf(codes.InvalidArgument, err.Error())
	}

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
		return response, status.Errorf(codes.InvalidArgument, err.Error())
	}

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
		return response, status.Errorf(codes.InvalidArgument, err.Error())
	}

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
		return response, status.Errorf(codes.InvalidArgument, err.Error())
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
		return response, status.Errorf(codes.InvalidArgument, err.Error())
	}

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
		return response, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return response, nil
}

// NewGrpcServer - creates new gRPC server instance
func NewGrpcServer(config *config.ServerConfig, storage storage.Repository) *GrpcServer {
	s := GrpcServer{
		config:  config,
		storage: storage,
	}
	return &s
}
