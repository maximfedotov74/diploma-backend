package model

type User struct {
	Id           int        `json:"id" example:"1" validate:"required"`
	Email        string     `json:"email" example:"makc-dgek@mail.ru" validate:"required"`
	PasswordHash string     `json:"-"`
	IsActivated  bool       `json:"is_activated" example:"false" validate:"required"`
	Patronymic   *string    `json:"patronymic"`
	FirstName    *string    `json:"first_name"`
	LastName     *string    `json:"last_name"`
	Roles        []UserRole `json:"roles" validate:"required"`
}

type UserRole struct {
	Id    int    `json:"id" example:"1" validate:"required"`
	Title string `json:"title" example:"User" validate:"required"`
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
