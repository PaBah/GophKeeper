package client

import (
	"context"
	"crypto/tls"
	"fmt"

	pb "github.com/PaBah/GophKeeper/internal/gen/proto/gophkeeper/v1"
	"github.com/PaBah/GophKeeper/internal/logger"
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
	serverAddress string
	conn          *grpc.ClientConn
	isAvailable   bool
}

func NewClientService(serverAddress string) *ClientService {
	return &ClientService{
		serverAddress: serverAddress,
	}
}

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

// IsNotAvailable checks if the server is not available.
func (c *ClientService) IsNotAvailable() bool {
	return !c.isAvailable
}

// IsAvailable checks if the server is available.
func (c *ClientService) IsAvailable() bool {
	return c.isAvailable
}

// TryToConnect attempts to establish a connection with the gRPC server.
// It sets up the connection and checks the server's availability.
func (c *ClientService) TryToConnect() bool {
	conn, err := grpc.Dial(c.serverAddress,
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
	})

	newCtx := metadata.NewOutgoingContext(ctx, md)

	return newCtx
}
