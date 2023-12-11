package auth

import (
	"fmt"

	"github.com/maximfedotov74/fiber-psql/internal/domain/session"
	"github.com/maximfedotov74/fiber-psql/internal/domain/user"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/jwt"
)

type UserService interface {
	GetUserById(id int) (*user.User, exception.Error)
	GetUserByEmail(email string) (*user.User, exception.Error)
	Create(dto user.CreateUserDto) (*user.UserCreatedResponse, exception.Error)
}
type SessionService interface {
	Parse(token string, tokenType jwt.TokenType) (*jwt.UserClaims, exception.Error)
	Sign(claims jwt.UserClaims) (*jwt.Tokens, exception.Error)
	CreateSession(dto session.CreateSessionDto) exception.Error
	FindSession(agent string, token string) (*session.Session, exception.Error)
}
type PasswordService interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashed string, pass string) bool
}

type MailService interface {
	SendActivationEmail(to string, subject string, link string) error
}

type AuthService struct {
	userService     UserService
	sessionService  SessionService
	passwordService PasswordService
	mailService     MailService
}

func NewAuthService(userService UserService, sessionService SessionService,
	passwordService PasswordService, mailService MailService) *AuthService {
	return &AuthService{
		userService:     userService,
		passwordService: passwordService,
		sessionService:  sessionService,
		mailService:     mailService,
	}
}

func (as *AuthService) Login(dto LoginDto, userAgent string) (*LoginResponse, exception.Error) {
	user, appErr := as.userService.GetUserByEmail(dto.Email)

	if appErr != nil {
		return nil, appErr
	}

	if isPasswordCorrect := as.passwordService.ComparePasswords(user.PasswordHash, dto.Password); !isPasswordCorrect {
		return nil, exception.NewErr(invalidCredentials, exception.STATUS_NOT_FOUND)
	}

	tokens, ex := as.sessionService.Sign(jwt.UserClaims{UserId: user.Id, UserAgent: userAgent})

	if ex != nil {
		return nil, ex
	}

	tokenDto := session.CreateSessionDto{UserId: user.Id, UserAgent: userAgent, Token: tokens.RefreshToken}
	appErr = as.sessionService.CreateSession(tokenDto)
	if appErr != nil {
		return nil, appErr
	}

	response := LoginResponse{Id: user.Id, Tokens: *tokens}

	return &response, nil
}

func (as *AuthService) Registration(dto user.CreateUserDto) (*int, exception.Error) {
	user, _ := as.userService.GetUserByEmail(dto.Email)

	if user != nil {
		return nil, exception.NewErr(userIsRegistered, exception.STATUS_BAD_REQUEST)
	}

	hash, err := as.passwordService.HashPassword(dto.Password)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}
	dto.Password = hash

	response, appErr := as.userService.Create(dto)

	if appErr != nil {
		return nil, appErr
	}

	link := fmt.Sprintf("/api/user/activate/%s", response.ActivationAccountLink)

	go as.mailService.SendActivationEmail(response.Email, "Активация аккаутна", link)

	return &response.Id, nil

}

func (as *AuthService) Refresh(refreshToken string, userAgent string) (*LoginResponse, exception.Error) {

	_, ex := as.sessionService.Parse(refreshToken, jwt.RefreshToken)

	if ex != nil {
		return nil, ex
	}

	dbToken, appErr := as.sessionService.FindSession(userAgent, refreshToken)

	if appErr != nil {
		return nil, appErr
	}

	user, appErr := as.userService.GetUserById(dbToken.UserId)
	if appErr != nil {
		return nil, appErr
	}

	claims := jwt.UserClaims{UserId: user.Id, UserAgent: dbToken.UserAgent}

	tokens, ex := as.sessionService.Sign(claims)

	if ex != nil {
		return nil, ex
	}

	tokenDto := session.CreateSessionDto{UserId: user.Id, UserAgent: dbToken.UserAgent, Token: tokens.RefreshToken}
	appErr = as.sessionService.CreateSession(tokenDto)

	if appErr != nil {
		return nil, appErr
	}

	response := LoginResponse{Id: user.Id, Tokens: *tokens}

	return &response, nil
}
