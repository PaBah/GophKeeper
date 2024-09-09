package models

import (
	"time"

	"github.com/PaBah/GophKeeper/internal/utils"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Credentials struct {
	ID          string    `json:"id"`
	ServiceName string    `json:"service_name"`
	Identity    string    `json:"identity"`
	Password    string    `json:"password"`
	UserID      string    `json:"-"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

type Card struct {
	ID             string    `json:"id"`
	Number         string    `json:"number"`
	ExpirationDate string    `json:"expiration_date"`
	HolderName     string    `json:"holder_name"`
	CVV            string    `json:"cvv"`
	UserID         string    `json:"-"`
	UploadedAt     time.Time `json:"uploaded_at"`
}

type File struct {
	Name       string    `json:"name"`
	Size       string    `json:"size"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func NewUser(email string, originalPassword string) User {
	return User{Email: email, Password: utils.PasswordHash(originalPassword)}
}

func NewCredentials(serviceName, identity, password string) Credentials {
	return Credentials{
		ServiceName: serviceName,
		Identity:    identity,
		Password:    password,
	}
}

func NewCard(number, expirationDate, holderName, cvv string) Card {
	return Card{
		Number:         number,
		ExpirationDate: expirationDate,
		HolderName:     holderName,
		CVV:            cvv,
	}
}
