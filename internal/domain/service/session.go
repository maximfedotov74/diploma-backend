package service

import (
	"context"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/jwt"
)

type sessionRepository interface {
	Create(ctx context.Context, dto model.CreateSessionDto) fall.Error
	FindByAgentAndToken(ctx context.Context, agent string, token string) (*model.Session, fall.Error)
	RemoveSession(ctx context.Context, userId int, sessionId int) fall.Error
	RemoveExceptCurrentSession(ctx context.Context, userId int, sessionId int) fall.Error
	RemoveAllSessions(ctx context.Context, userId int) fall.Error
	GetUserSessions(ctx context.Context, userId int, agent string) (*model.UserSessionsResponse, fall.Error)
	FindByAgentAndUserId(ctx context.Context, agent string, userId int) (*model.Session, fall.Error)
	RemoveSessionByToken(ctx context.Context, token string) fall.Error
}

type jwtService interface {
	Sign(claims jwt.UserClaims) (jwt.Tokens, error)
	Parse(token string, tokenType jwt.TokenType) (*jwt.UserClaims, error)
}

type SessionService struct {
	repo sessionRepository
	jwt  jwtService
}

func NewSessionService(repo sessionRepository, jwt jwtService) *SessionService {
	return &SessionService{repo: repo, jwt: jwt}
}

func (s *SessionService) FindByAgentAndUserId(ctx context.Context, agent string, userId int) (*model.Session, fall.Error) {
	return s.repo.FindByAgentAndUserId(ctx, agent, userId)

}

func (s *SessionService) GetUserSessions(ctx context.Context, userId int, token string) (*model.UserSessionsResponse, fall.Error) {
	return s.repo.GetUserSessions(ctx, userId, token)
}

func (s *SessionService) Create(ctx context.Context, dto model.CreateSessionDto) fall.Error {
	return s.repo.Create(ctx, dto)
}

func (s *SessionService) FindByAgentAndToken(ctx context.Context, agent string, token string) (*model.Session, fall.Error) {
	return s.repo.FindByAgentAndToken(ctx, agent, token)
}

func (s *SessionService) RemoveSessionByToken(ctx context.Context, token string) fall.Error {
	return s.repo.RemoveSessionByToken(ctx, token)
}

func (s *SessionService) RemoveSession(ctx context.Context, userId int, sessionId int) fall.Error {
	return s.repo.RemoveSession(ctx, userId, sessionId)
}

func (s *SessionService) RemoveExceptCurrentSession(ctx context.Context, userId int, sessionId int) fall.Error {
	return s.repo.RemoveExceptCurrentSession(ctx, userId, sessionId)
}
func (s *SessionService) RemoveAllSessions(ctx context.Context, userId int) fall.Error {
	return s.repo.RemoveAllSessions(ctx, userId)
}

func (s *SessionService) Sign(claims jwt.UserClaims) (*jwt.Tokens, fall.Error) {

	tokens, err := s.jwt.Sign(claims)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	return &tokens, nil
}

func (s *SessionService) Parse(token string, tokenType jwt.TokenType) (*jwt.UserClaims, fall.Error) {
	claims, err := s.jwt.Parse(token, tokenType)
	if err != nil {
		return nil, fall.NewErr(err.Error(), fall.STATUS_UNAUTHORIZED)
	}
	return claims, nil
}
