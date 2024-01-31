package model

type RoleUser struct {
	Id         *int    `json:"user_id" validate:"required" example:"1"`
	Email      *string `json:"email" validate:"required" example:"example@mail.ru"`
	UserRoleId *int    `json:"-"`
	RoleId     *int    `json:"-"`
}

type Role struct {
	Id    int        `json:"role_id" db:"role_id" example:"1" validate:"required"`
	Title string     `json:"title" db:"title" example:"ADMIN" validate:"required"`
	Users []RoleUser `json:"users" validate:"required"`
}

type AddRoleToUserDto struct {
	Title  string `json:"title" validate:"required,min=3" example:"ADMIN"`
	UserId int    `json:"user_id" validate:"required,min=1" example:"1"`
}

type CreateRoleDto struct {
	Title string `json:"title" validate:"required,min=6,max=55" db:"title" example:"ADMIN"`
}
