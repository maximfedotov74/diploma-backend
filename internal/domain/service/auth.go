package service

import (
	"context"
	"fmt"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/jwt"
	"github.com/maximfedotov74/diploma-backend/internal/shared/password"
)

type authUserService interface {
	Create(ctx context.Context, dto model.CreateUserDto) (*model.CreatedUserResponse, fall.Error)
	FindByEmail(ctx context.Context, email string) (*model.User, fall.Error)
	FindById(ctx context.Context, id int) (*model.User, fall.Error)
}
type authSessionService interface {
	Parse(token string, tokenType jwt.TokenType) (*jwt.UserClaims, fall.Error)
	Sign(claims jwt.UserClaims) (*jwt.Tokens, fall.Error)
	Create(ctx context.Context, dto model.CreateSessionDto) fall.Error
	FindByAgentAndToken(ctx context.Context, agent string, token string) (*model.Session, fall.Error)
	RemoveSessionByToken(ctx context.Context, token string) fall.Error
}

type authMailService interface {
	SendActivationEmail(to string, subject string, link string) error
}

type AuthService struct {
	userService    authUserService
	sessionService authSessionService
	mailService    authMailService
}

func NewAuthService(userService authUserService, sessionService authSessionService,
	mailService authMailService) *AuthService {
	return &AuthService{
		userService:    userService,
		sessionService: sessionService,
		mailService:    mailService,
	}
}

func (as *AuthService) Login(ctx context.Context, dto model.LoginDto, userAgent string) (*model.LoginResponse, fall.Error) {
	user, appErr := as.userService.FindByEmail(ctx, dto.Email)

	if appErr != nil {
		return nil, appErr
	}

	if isPasswordCorrect := password.ComparePasswords(user.PasswordHash, dto.Password); !isPasswordCorrect {
		return nil, fall.NewErr(msg.InvalidCredentials, fall.STATUS_NOT_FOUND)
	}

	tokens, ex := as.sessionService.Sign(jwt.UserClaims{UserId: user.Id, UserAgent: userAgent})

	if ex != nil {
		return nil, ex
	}

	tokenDto := model.CreateSessionDto{UserId: user.Id, UserAgent: userAgent, Token: tokens.RefreshToken}
	appErr = as.sessionService.Create(ctx, tokenDto)
	if appErr != nil {
		return nil, appErr
	}

	response := model.LoginResponse{Id: user.Id, Tokens: *tokens}

	return &response, nil
}

func (as *AuthService) Registration(ctx context.Context, dto model.CreateUserDto) (*int, fall.Error) {
	user, _ := as.userService.FindByEmail(ctx, dto.Email)

	if user != nil {
		return nil, fall.NewErr(msg.UserIsRegistered, fall.STATUS_BAD_REQUEST)
	}

	hash, err := password.HashPassword(dto.Password)

	if err != nil {
		return nil, fall.NewErr(err.Error(), 500)
	}
	dto.Password = hash

	response, appErr := as.userService.Create(ctx, dto)

	if appErr != nil {
		return nil, appErr
	}

	link := fmt.Sprintf("/api/user/activate/%s", response.Link)

	go as.mailService.SendActivationEmail(response.Email, "Активация аккаутна", link)

	return &response.Id, nil

}

func (as *AuthService) Refresh(ctx context.Context, refreshToken string) (*model.LoginResponse, fall.Error) {

	parsed, ex := as.sessionService.Parse(refreshToken, jwt.RefreshToken)

	if ex != nil {
		return nil, ex
	}

	dbToken, appErr := as.sessionService.FindByAgentAndToken(ctx, parsed.UserAgent, refreshToken)

	if appErr != nil {
		return nil, appErr
	}

	user, appErr := as.userService.FindById(ctx, dbToken.UserId)
	if appErr != nil {
		return nil, appErr
	}

	claims := jwt.UserClaims{UserId: user.Id, UserAgent: dbToken.UserAgent}

	tokens, ex := as.sessionService.Sign(claims)

	if ex != nil {
		return nil, ex
	}

	tokenDto := model.CreateSessionDto{UserId: user.Id, UserAgent: dbToken.UserAgent, Token: tokens.RefreshToken}
	appErr = as.sessionService.Create(ctx, tokenDto)

	if appErr != nil {
		return nil, appErr
	}

	response := model.LoginResponse{Id: user.Id, Tokens: *tokens}

	return &response, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) fall.Error {
	return s.sessionService.RemoveSessionByToken(ctx, token)
}
