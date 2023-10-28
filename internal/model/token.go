package model

type Token struct {
	TokenId   int    `json:"token_id" db:"token_id"`
	UserId    int    `json:"user_id" db:"user_id"`
	UserAgent string `json:"user_agent" db:"user_agent"`
	Token     string `json:"token" db:"token"`
}

type CreateToken struct {
	UserId    int    `json:"user_id" db:"user_id"`
	UserAgent string `json:"user_agent" db:"user_agent"`
	Token     string `json:"token" db:"token"`
}
