package model

import (
	"github.com/maximfedotov74/fiber-psql/pkg/token"
)

type User struct {
	Id           int    `json:"user_id" db:"user_id" example:"1"`
	Email        string `json:"email" validate:"required,email" db:"email" example:"makc@mail.ru"`
	PasswordHash string `json:"password_hash" validate:"required,min=6,max=100" db:"password_hash" example:"sdfsdfs222"`
	IsActivated  bool   `json:"is_activated" db:"is_activated" example:"false"`
	Roles        []Role `json:"roles"`
}

type RegistrationResponse struct {
	Id int `json:"id" example:"1"`
}

type CreateUserDto struct {
	Email    string `json:"email" validate:"required,email" example:"makc@mail.ru"`
	Password string `json:"password" validate:"required,min=6,max=100" example:"sdfsdfs222"`
}

type LoginDto struct {
	Email    string `json:"email" validate:"required,email" example:"makc@mail.ru"`
	Password string `json:"password" validate:"required,min=6,max=100" example:"sdfsdfs222"`
}

type LoginResponse struct {
	Id     int          `json:"user_id" db:"user_id" example:"1"`
	Tokens token.Tokens `json:"tokens"`
	Roles  []Role       `json:"roles"`
}
