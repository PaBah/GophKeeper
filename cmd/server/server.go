package main

import (
	"context"
	"errors"
	"fmt"

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

// SignIn - handler for shortening URL
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

// SignUp - handler for shortening URL
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

// NewGrpcServer - creates new gRPC server instance
func NewGrpcServer(config *config.ServerConfig, storage storage.Repository) *GrpcServer {
	s := GrpcServer{
		config:  config,
		storage: storage,
	}
	return &s
}
