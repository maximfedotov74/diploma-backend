package service

import (
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
		return lib.NewErr(messages.TOKEN_CREATE_ERROR, 500)
	}

	return nil
}

func (ts *TokenService) FindToken(agent string, token string) (*model.Token, lib.Error) {

	dbToken, err := ts.repo.FindByAgentAndToken(agent, token)

	if err != nil {
		return nil, lib.NewErr(messages.TOKEN_NOT_FOUND, 404)
	}
	return dbToken, nil

}

func (ts *TokenService) RemoveToken() error {
	return nil
}

func (ts *TokenService) Refresh(refreshToken string) error {

	_, err := ts.jwtService.Parse(refreshToken, token.RefreshToken)

	if err != nil {
		return err
	}

	return nil
}

func (ts *TokenService) Sign(claims token.UserClaims) (token.Tokens, error) {
	return ts.jwtService.Sign(claims)
}

func (ts *TokenService) Parse(token string, tokenType token.TokenType) (*token.UserClaims, error) {
	return ts.jwtService.Parse(token, tokenType)
}
