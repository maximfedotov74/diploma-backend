package model

import "github.com/go-playground/validator/v10"

type UserGender string

const (
	Men   UserGender = "men"
	Women UserGender = "women"
)

func UserGenderEnumValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	switch value {
	case string(Men), string(Women):
		return true
	}
	return false
}

type GetAllUsersResponse struct {
	Users []*User `json:"users" validate:"required"`
	Total int     `json:"total" validate:"required"`
}

type User struct {
	Id           int         `json:"id" example:"1" validate:"required"`
	Email        string      `json:"email" example:"makc-dgek@mail.ru" validate:"required"`
	PasswordHash string      `json:"-"`
	IsActivated  bool        `json:"is_activated" example:"false" validate:"required"`
	Patronymic   *string     `json:"patronymic" validate:"omitempty,min=3"`
	FirstName    *string     `json:"first_name" validate:"omitempty,min=1"`
	LastName     *string     `json:"last_name" validate:"omitempty,min=1"`
	Roles        []UserRole  `json:"roles" validate:"required"`
	Gender       *UserGender `json:"gender"`
	AvatarPath   *string     `json:"avatar_path" validate:"omitempty,filepath"`
}

type UserRole struct {
	Id         *int    `json:"id" example:"1" validate:"required"`
	Title      *string `json:"title" example:"User" validate:"required"`
	UserId     *int    `json:"-"`
	UserRoleId *int    `json:"-"`
}

type CreatedUserResponse struct {
	Id    int    `json:"id" example:"1" validate:"required"`
	Email string `json:"email" example:"makc-dgek@mail.ru" validate:"required"`
	Link  string `json:"link"`
}

type CreateUserDto struct {
	Email    string `json:"email" example:"makc-dgek@mail.ru" validate:"required,email"`
	Password string `json:"password" example:"1234567890" validate:"required,min=6"`
}

type UpdateUserDto struct {
	AvatarPath *string     `json:"avatar_path" validate:"omitempty,filepath"`
	Gender     *UserGender `json:"gender" validate:"omitempty,userGenderEnumValidation"`
	Patronymic *string     `json:"patronymic" validate:"omitempty,min=3"`
	FirstName  *string     `json:"first_name" validate:"omitempty,min=1"`
	LastName   *string     `json:"last_name" validate:"omitempty,min=3"`
}

type ChangePasswordCode struct {
	ChangePasswordCodeId int    `json:"change_password_code_id" db:"change_password_code_id" validate:"required"`
	Code                 string `json:"code" db:"code" validate:"required"`
	UserId               int    `json:"user_id" db:"user_id" validate:"required"`
}

type ConfirmChangePasswordDto struct {
	Code string `json:"code" validate:"required,min=6,max=6" example:"123456"`
}

type ChangePasswordDto struct {
	OldPassword string `json:"old_password" validate:"required,min=6,max=100" example:"sdfsdfs222"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=100" example:"sdfsdfs222"`
}
