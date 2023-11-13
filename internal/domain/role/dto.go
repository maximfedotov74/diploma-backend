package role

type AddRoleToUserDto struct {
	Title  string `json:"title" validate:"required,min=3" example:"ADMIN"`
	UserId int    `json:"user_id" validate:"required,min=1" example:"1"`
}

type CreateRoleDto struct {
	Title string `json:"title" validate:"required,min=6,max=55" db:"title" example:"ADMIN"`
}
