package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/maximfedotov74/fiber-psql/internal/cfg"
)

type Token interface {
	Sign(id int) (Tokens, error)
	Parse(token string, tokenType TokenType) (int, error)
}

type Claims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

type TokenService struct {
	config *cfg.Config
}

type Tokens struct {
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	AccessExpTime  time.Time `json:"-"`
	RefreshExpTime time.Time `json:"-"`
}

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)

func New(cfg *cfg.Config) *TokenService {
	return &TokenService{
		config: cfg,
	}
}

func (ts *TokenService) Sign(id int) (Tokens, error) {

	var tokens Tokens

	var accessSecret = ts.config.AccessTokenSecret
	var accessExpTime time.Time = time.Now().Add(time.Minute * time.Duration(ts.config.AccessTokenExp))

	var refreshExpTime time.Time = time.Now().AddDate(0, 0, ts.config.RefreshTokenExp)
	var refreshSecret string = ts.config.RefreshTokenSecret

	accessClaims := Claims{
		UserId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshClaims := Claims{
		UserId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	acccessTokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshTokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessToken, err := acccessTokenObject.SignedString([]byte(accessSecret))
	if err != nil {
		return tokens, err
	}

	refreshToken, err := refreshTokenObject.SignedString([]byte(refreshSecret))
	if err != nil {
		return tokens, err
	}

	tokens = Tokens{AccessToken: accessToken, RefreshToken: refreshToken, AccessExpTime: accessExpTime, RefreshExpTime: refreshExpTime}
	return tokens, nil

}

func (ts *TokenService) Parse(token string, tokenType TokenType) (int, error) {
	result, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid token!!")
		}
		var secret string
		if tokenType == AccessToken {
			secret = ts.config.AccessTokenSecret
		} else {
			secret = ts.config.RefreshTokenSecret
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	if !result.Valid {
		return 0, errors.New("Token is not valid!")
	}

	claims, ok := result.Claims.(*Claims)

	if !ok {
		return 0, errors.New("Cannot parse claims to struct!")
	}

	return claims.UserId, nil

}
