package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
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

func (c *ClientService) CreateCredentials(ctx context.Context, serviceName, identity, password string) error {
	resp, err := c.client.CreateCredentials(c.getCtx(ctx, c.token), &pb.CreateCredentialsRequest{
		ServiceName: serviceName,
		Identity:    identity,
		Password:    password,
	})
	if err != nil {
		return fmt.Errorf("CreateCredentials: %w", err)
	}
	log.Println("CreateCredentials", resp.Id)
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
	qwe, err := c.client.CreateCard(c.getCtx(ctx, c.token), &pb.CreateCardRequest{
		Number:         number,
		ExpirationDate: expirationDate,
		HolderName:     holderName,
		Cvv:            cvv,
	})

	log.Println("CreateCard", qwe)
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
