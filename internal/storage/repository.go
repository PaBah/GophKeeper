package storage

import (
	"errors"
)

// ErrConflict - error when user tries to save already existing data
var ErrConflict = errors.New("data conflict")

// Repository - interface over Repository pattern for system storage
type Repository interface {
}
