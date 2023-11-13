package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	AccessExpTime  time.Time `json:"-"`
	RefreshExpTime time.Time `json:"-"`
}

type TokenType int

type UserClaims struct {
	UserId    int    `json:"user_id"`
	UserAgent string `json:"user_agent"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserClaims
}
