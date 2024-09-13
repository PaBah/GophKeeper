package storage

import (
	"context"
	"errors"

	"github.com/PaBah/GophKeeper/internal/models"
)

// ErrAlreadyExists - error when user tries to save already existing data
var ErrAlreadyExists = errors.New("already exists")

// Repository - interface over Repository pattern for system storage
type Repository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	AuthorizeUser(ctx context.Context, email string) (models.User, error)
	CreateCredentials(ctx context.Context, credentials models.Credentials) (models.Credentials, error)
	GetCredentials(ctx context.Context) ([]models.Credentials, error)
	UpdateCredentials(ctx context.Context, credentials models.Credentials) (models.Credentials, error)
	DeleteCredentials(ctx context.Context, credentialsID string) error

	CreateCard(ctx context.Context, card models.Card) (models.Card, error)
	GetCards(ctx context.Context) ([]models.Card, error)
	UpdateCard(ctx context.Context, card models.Card) (models.Card, error)
	DeleteCard(ctx context.Context, cardID string) error
}
