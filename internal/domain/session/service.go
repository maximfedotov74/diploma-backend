package session

import (
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/jwt"
)

// todo add remove sessions

type JwtTokenService interface{}

type Repository interface {
	CreateSession(CreateSessionDto) exception.Error
	FindByAgentAndToken(agent string, token string) (*Session, exception.Error)
	RemoveSession(token string, agent string) exception.Error
	RemoveExceptCurrentSession(userId int, agent string) exception.Error
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
		return err
	}

	return nil
}

func (ts *TokenService) FindSession(agent string, token string) (*Session, exception.Error) {

	dbToken, err := ts.repo.FindByAgentAndToken(agent, token)

	if err != nil {
		return nil, err
	}
	return dbToken, nil

}

func (ts *TokenService) RemoveSession(token string, agent string) exception.Error {
	err := ts.repo.RemoveSession(token, agent)
	return err
}

func (ts *TokenService) RemoveExceptCurrentSession(userId int, agent string) exception.Error {
	err := ts.repo.RemoveExceptCurrentSession(userId, agent)
	return err
}

func (ts *TokenService) Refresh(refreshToken string) error {

	_, err := ts.jwtService.Parse(refreshToken, jwt.RefreshToken)

	if err != nil {
		return err
	}

	return nil
}

func (ts *TokenService) Sign(claims jwt.UserClaims) (*jwt.Tokens, exception.Error) {

	tokens, err := ts.jwtService.Sign(claims)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}
	return &tokens, nil
}

func (ts *TokenService) Parse(token string, tokenType jwt.TokenType) (*jwt.UserClaims, exception.Error) {
	claims, err := ts.jwtService.Parse(token, tokenType)
	if err != nil {
		return nil, exception.NewErr(err.Error(), exception.STATUS_UNAUTHORIZED)
	}
	return claims, nil
}
