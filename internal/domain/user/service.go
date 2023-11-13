package user

import (
	"github.com/maximfedotov74/fiber-psql/internal/domain/session"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/jwt"
	"github.com/maximfedotov74/fiber-psql/internal/shared/messages"
	"github.com/maximfedotov74/fiber-psql/internal/shared/models"
)

type MailService interface {
	SendChangePasswordEmail(to string, subject string, code string) error
}

type SessionServie interface {
	CreateSession(dto session.CreateSessionDto) exception.Error
	Sign(claims jwt.UserClaims) (jwt.Tokens, error)
}

type Repository interface {
	Create(password string, email string) (*UserCreatedResponse, error)
	FindActivationLink(link string) (*int, error)
	ActivateUser(id *int) error
	GetUserById(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	ChangePassword(userId int, newPassword string) error
	FindChangePasswordCode(userId int, code string) (*ChangePasswordCode, error)
	CreateChangePasswordCode(userId int) (*string, error)
}

type PasswordService interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashed string, pass string) bool
}

type UserService struct {
	repo            Repository
	sessionService  SessionServie
	passwordService PasswordService
	mailService     MailService
}

func NewUserService(repo Repository, sessionService SessionServie, mailService MailService, passwordService PasswordService) *UserService {
	return &UserService{
		repo:            repo,
		sessionService:  sessionService,
		mailService:     mailService,
		passwordService: passwordService,
	}
}

func (us *UserService) GetAll() {}

func (us *UserService) Create(dto CreateUserDto) (*UserCreatedResponse, exception.Error) {

	response, err := us.repo.Create(dto.Password, dto.Email)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}

	return response, nil

}

func (us *UserService) Activate(activationLink string) exception.Error {
	id, err := us.repo.FindActivationLink(activationLink)

	if err != nil {
		return exception.NewErr(err.Error(), 404)
	}

	err = us.repo.ActivateUser(id)
	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}

	return nil
}

func (us *UserService) GetUserById(id int) (*User, exception.Error) {
	user, err := us.repo.GetUserById(id)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}

	if user == nil {
		return nil, exception.NewErr(messages.USER_NOT_FOUND, 404)
	}

	return user, nil
}

func (us *UserService) GetUserByEmail(email string) (*User, exception.Error) {
	user, err := us.repo.GetUserByEmail(email)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}

	if user == nil {
		return nil, exception.NewErr(messages.USER_NOT_FOUND, 404)
	}

	return user, nil
}

func (us *UserService) CreateChangePasswordCode(userId int) exception.Error {

	currentUser, appErr := us.GetUserById(userId)
	if appErr != nil {
		return appErr
	}

	code, err := us.repo.CreateChangePasswordCode(currentUser.Id)

	if err != nil {
		return exception.NewErr(err.Error(), 500)
	}

	go us.mailService.SendChangePasswordEmail(currentUser.Email, "Код для смены пароля", *code)

	return nil
}

func (us *UserService) ChangePassword(dto ChangePasswordDto, contextData *models.UserContextData) (*jwt.Tokens, exception.Error) {

	user, appErr := us.GetUserById(contextData.UserId)

	if appErr != nil {
		return nil, appErr
	}

	oldMatch := us.passwordService.ComparePasswords(user.PasswordHash, dto.OldPassword)

	if !oldMatch {
		return nil, exception.NewErr(messages.BAD_PASSWORD, 400)
	}

	newMatch := us.passwordService.ComparePasswords(user.PasswordHash, dto.NewPassword)

	if newMatch {
		return nil, exception.NewErr(messages.BAD_NEW_PASSWORD, 400)

	}

	newHash, err := us.passwordService.HashPassword(dto.NewPassword)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}

	_, err = us.repo.FindChangePasswordCode(user.Id, dto.Code)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 404)
	}

	err = us.repo.ChangePassword(user.Id, newHash)

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}

	tokens, err := us.sessionService.Sign(jwt.UserClaims{UserId: user.Id, UserAgent: contextData.UserAgent})

	if err != nil {
		return nil, exception.NewErr(err.Error(), 500)
	}

	tokenDto := session.CreateSessionDto{UserId: user.Id, UserAgent: contextData.UserAgent, Token: tokens.RefreshToken}
	appErr = us.sessionService.CreateSession(tokenDto)

	if appErr != nil {
		return nil, appErr
	}

	return &tokens, nil
}
