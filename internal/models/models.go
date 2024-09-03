package models

import "github.com/PaBah/GophKeeper/internal/utils"

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func NewUser(email string, originalPassword string) User {
	return User{Email: email, Password: utils.PasswordHash(originalPassword)}
}
