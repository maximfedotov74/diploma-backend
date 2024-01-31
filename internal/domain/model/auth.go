package model

import "github.com/maximfedotov74/diploma-backend/internal/shared/jwt"

type RegistrationResponse struct {
	Id int `json:"id" example:"1" validate:"required"`
}

type LoginResponse struct {
	Id     int        `json:"user_id" db:"user_id" example:"1" validate:"required"`
	Tokens jwt.Tokens `json:"tokens" validate:"required"`
}

type LoginDto struct {
	Email    string `json:"email" validate:"required,email" example:"makc@mail.ru"`
	Password string `json:"password" validate:"required,min=6,max=100" example:"sdfsdfs222"`
}
