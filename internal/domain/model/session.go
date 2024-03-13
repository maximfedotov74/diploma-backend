package model

import "time"

type Session struct {
	SessionId int       `json:"session_id" db:"session_id" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
	UserId    int       `json:"user_id" db:"user_id" validate:"required"`
	UserAgent string    `json:"user_agent" db:"user_agent" validate:"required"`
	Token     string    `json:"-" db:"token" validate:"required"`
}

type CreateSessionDto struct {
	UserId    int    `json:"user_id" db:"user_id"`
	UserAgent string `json:"user_agent" db:"user_agent"`
	Token     string `json:"token" db:"token"`
}

type LocalSession struct {
	UserId    int        `json:"user_id"`
	UserAgent string     `json:"user_agent"`
	Email     string     `json:"email"`
	Roles     []UserRole `json:"roles"`
}

type UserSessionsResponse struct {
	Current *Session  `json:"current" validate:"required"`
	All     []Session `json:"sessions" validate:"required"`
}
