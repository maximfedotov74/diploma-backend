package model

type CreateRoleDto struct {
	Title string `json:"title" validate:"required,min=6,max=55" db:"title"`
}

type Role struct {
	Id    int    `json:"role_id" db:"role_id" example:"1"`
	Title string `json:"title" db:"title" example:"ADMIN"`
}

type AddRoleToUserDto struct {
	Title  string `json:"title" validate:"required,min=3"`
	UserId int    `json:"user_id" validate:"required,min=1"`
}
