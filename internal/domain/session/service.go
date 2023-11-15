package session

import (
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/jwt"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
)

// todo add remove sessions

type JwtTokenService interface{}

type Repository interface {
	CreateSession(CreateSessionDto) error
	FindByAgentAndToken(agent string, token string) (*Session, error)
}

type JwtService interface {
	Parse(token string, tokenType jwt.TokenType) (*jwt.UserClaims, error)
	Sign(claims jwt.UserClaims) (jwt.Tokens, error)
}

type TokenService struct {
	jwtService JwtService
	repo       Repository
}

func NewSessionService(repo Repository, jwt JwtService) *TokenService {

	return &TokenService{
		jwtService: jwt,
		repo:       repo,
	}
}

func (ts *TokenService) CreateSession(dto CreateSessionDto) exception.Error {
	err := ts.repo.CreateSession(dto)

	if err != nil {
		return exception.NewErr(messages.TOKEN_CREATE_ERROR, 500)
	}

	return nil
}

func (ts *TokenService) FindSession(agent string, token string) (*Session, exception.Error) {

	dbToken, err := ts.repo.FindByAgentAndToken(agent, token)

	if err != nil {
		return nil, exception.NewErr(messages.TOKEN_NOT_FOUND, 404)
	}
	return dbToken, nil

}

func (ts *TokenService) RemoveSession() error {
	return nil
}

func (ts *TokenService) Refresh(refreshToken string) error {

	_, err := ts.jwtService.Parse(refreshToken, jwt.RefreshToken)

	if err != nil {
		return err
	}

	return nil
}

func (ts *TokenService) Sign(claims jwt.UserClaims) (jwt.Tokens, error) {
	return ts.jwtService.Sign(claims)
}

func (ts *TokenService) Parse(token string, tokenType jwt.TokenType) (*jwt.UserClaims, error) {
	return ts.jwtService.Parse(token, tokenType)
}
