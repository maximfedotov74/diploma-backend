package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtConfig struct {
	RefreshTokenExp    int
	AccessTokenExp     int
	RefreshTokenSecret string
	AccessTokenSecret  string
}

type JwtService struct {
	config JwtConfig
}

const tokenInvalid = "Ошибка при валидации токена!"
const parseClaimsError = "Данные записанные в токен не соответствуют требуемым!"

const (
	AccessToken TokenType = iota
	RefreshToken
)

func NewJwtService(cfg JwtConfig) *JwtService {
	return &JwtService{
		config: cfg,
	}
}

func (ts *JwtService) Sign(claims UserClaims) (Tokens, error) {

	var tokens Tokens

	var accessSecret = ts.config.AccessTokenSecret
	var accessExpTime time.Time = time.Now().Add(time.Minute * time.Duration(ts.config.AccessTokenExp))

	var refreshExpTime time.Time = time.Now().AddDate(0, 0, ts.config.RefreshTokenExp)
	var refreshSecret string = ts.config.RefreshTokenSecret

	accessClaims := Claims{
		UserClaims: claims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshClaims := Claims{
		UserClaims: claims,
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

func (ts *JwtService) Parse(token string, tokenType TokenType) (*UserClaims, error) {
	result, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(tokenInvalid)
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
		return nil, err
	}

	if !result.Valid {
		return nil, errors.New(tokenInvalid)
	}

	claims, ok := result.Claims.(*Claims)

	if !ok {
		return nil, errors.New(parseClaimsError)
	}

	return &claims.UserClaims, nil

}
