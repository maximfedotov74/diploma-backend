package user

import (
	"github.com/maximfedotov74/fiber-psql/internal/domain/role"
)

type User struct {
	Id           int         `json:"user_id" db:"user_id" example:"1"`
	Email        string      `json:"email" validate:"required,email" db:"email" example:"makc@mail.ru"`
	PasswordHash string      `json:"-" validate:"required,min=6,max=100" db:"password_hash" example:"sdfsdfs222"`
	IsActivated  bool        `json:"is_activated" db:"is_activated" example:"false"`
	Roles        []role.Role `json:"roles"`
}

type UserCreatedResponse struct {
	Id                    int    `json:"id"`
	ActivationAccountLink string `json:"activation_account_link"`
	Email                 string `json:"email"`
}

type UserSettings struct {
	Id                    int     `json:"user_settings_id" db:"user_settings_id" example:"1"`
	ActivationAccountLink *string `json:"activation_account_link" db:"activation_account_link"`
	UserId                int     `json:"user_id" db:"user_id"`
}

type ChangePasswordCode struct {
	ChangePasswordCodeId int    `json:"change_password_code_id" db:"change_password_code_id"`
	Code                 string `json:"code" db:"code"`
	UserId               int    `json:"user_id" db:"user_id"`
}
