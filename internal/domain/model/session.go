package model

type Session struct {
	SessionId int    `json:"session_id" db:"session_id" validate:"required"`
	UserId    int    `json:"user_id" db:"user_id" validate:"required"`
	UserAgent string `json:"user_agent" db:"user_agent" validate:"required"`
	Token     string `json:"token" db:"token" validate:"required"`
}

type CreateSessionDto struct {
	UserId    int    `json:"user_id" db:"user_id"`
	UserAgent string `json:"user_agent" db:"user_agent"`
	Token     string `json:"token" db:"token"`
}

type LocalSession struct {
	UserId    int
	UserAgent string
	Roles     []UserRole
}
