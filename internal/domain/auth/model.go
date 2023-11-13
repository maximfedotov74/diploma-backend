package auth

import "github.com/maximfedotov74/fiber-psql/internal/shared/jwt"

type RegistrationResponse struct {
	Id int `json:"id" example:"1"`
}

type LoginResponse struct {
	Id     int        `json:"user_id" db:"user_id" example:"1"`
	Tokens jwt.Tokens `json:"tokens"`
}
