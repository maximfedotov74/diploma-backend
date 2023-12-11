package user

import (
	"github.com/maximfedotov74/fiber-psql/internal/cfg"
	"github.com/maximfedotov74/fiber-psql/internal/domain/session"
	exception "github.com/maximfedotov74/fiber-psql/internal/shared/error"
	"github.com/maximfedotov74/fiber-psql/internal/shared/jwt"
	"github.com/maximfedotov74/fiber-psql/internal/shared/models"
)

type MailService interface {
	SendChangePasswordEmail(to string, subject string, code string) error
}

type SessionServie interface {
	CreateSession(dto session.CreateSessionDto) exception.Error
	Sign(claims jwt.UserClaims) (*jwt.Tokens, exception.Error)
}

type Repository interface {
	Create(password string, email string) (*UserCreatedResponse, exception.Error)
	FindActivationLink(link string) (*int, exception.Error)
	ActivateUser(id *int) exception.Error
	GetUserById(id int) (*User, exception.Error)
	GetUserByEmail(email string) (*User, exception.Error)
	ChangePassword(userId int, newPassword string) exception.Error
	FindChangePasswordCode(userId int, code string) (*ChangePasswordCode, exception.Error)
	CreateChangePasswordCode(userId int) (*string, exception.Error)
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
	cfgService      *cfg.Config
}

func NewUserService(repo Repository, sessionService SessionServie, mailService MailService, passwordService PasswordService,
) *UserService {
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
		return nil, err
	}

	return response, nil

}

func (us *UserService) Activate(activationLink string) exception.Error {
	id, err := us.repo.FindActivationLink(activationLink)

	if err != nil {
		return err
	}

	err = us.repo.ActivateUser(id)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) GetUserById(id int) (*User, exception.Error) {

	user, err := us.repo.GetUserById(id)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (us *UserService) GetUserByEmail(email string) (*User, exception.Error) {
	user, err := us.repo.GetUserByEmail(email)

	if err != nil {
		return nil, err
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
		return err
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
		return nil, exception.NewErr(badPassword, exception.STATUS_BAD_REQUEST)
	}

	newMatch := us.passwordService.ComparePasswords(user.PasswordHash, dto.NewPassword)

	if newMatch {
		return nil, exception.NewErr(badNewPassword, exception.STATUS_BAD_REQUEST)

	}

	newHash, err := us.passwordService.HashPassword(dto.NewPassword)

	if err != nil {
		return nil, exception.ServerError(err.Error())
	}

	_, ex := us.repo.FindChangePasswordCode(user.Id, dto.Code)

	if ex != nil {
		return nil, ex
	}

	ex = us.repo.ChangePassword(user.Id, newHash)

	if ex != nil {
		return nil, ex
	}

	tokens, ex := us.sessionService.Sign(jwt.UserClaims{UserId: user.Id, UserAgent: contextData.UserAgent})

	if ex != nil {
		return nil, ex
	}

	tokenDto := session.CreateSessionDto{UserId: user.Id, UserAgent: contextData.UserAgent, Token: tokens.RefreshToken}
	appErr = us.sessionService.CreateSession(tokenDto)

	if appErr != nil {
		return nil, appErr
	}

	return tokens, nil
}
