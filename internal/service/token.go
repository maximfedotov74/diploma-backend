package service

import (
	"log"

	"github.com/maximfedotov74/fiber-psql/internal/cfg"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/repository"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
	"github.com/maximfedotov74/fiber-psql/pkg/token"
)

type TokenService struct {
	jwtService *token.TokenService
	repo       repository.Token
}

func NewTokenService(config *cfg.Config, repo repository.Token) *TokenService {
	jwtTokenService := token.New(config)

	return &TokenService{
		jwtService: jwtTokenService,
		repo:       repo,
	}
}

func (ts *TokenService) Create(dto model.CreateToken) lib.Error {
	err := ts.repo.CreateToken(dto)

	if err != nil {
		log.Println(err)
		return lib.NewErr(messages.TOKEN_CREATE_ERROR, 500)
	}

	return nil
}

func (ts *TokenService) FindToken() error {
	return nil
}

func (ts *TokenService) RemoveToken() error {
	return nil
}

func (ts *TokenService) Refresh() error {
	return nil
}

func (ts *TokenService) Sign(claims token.UserClaims) (token.Tokens, error) {
	return ts.jwtService.Sign(claims)
}

func (ts *TokenService) Parse(token string, tokenType token.TokenType) (*token.UserClaims, error) {
	return ts.jwtService.Parse(token, tokenType)
}
