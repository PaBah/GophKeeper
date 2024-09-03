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
}
