package model

type СompletedOperation struct {
	Completed bool `json:"completed" example:"true"`
}

type UserContextData struct {
	UserId    int    `json:"user_id" db:"user_id" example:"1"`
	Roles     []Role `json:"roles"`
	UserAgent string `json:"user_agent"`
}
