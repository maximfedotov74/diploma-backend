package models

type role struct {
	Id    int    `json:"role_id" db:"role_id" example:"1"`
	Title string `json:"title" db:"title" example:"ADMIN"`
}

type UserContextData struct {
	UserId    int    `json:"user_id"`
	UserAgent string `json:"user_agent"`
	Roles     []role
}
