package service

import (
	"context"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
	"github.com/maximfedotov74/diploma-backend/internal/shared/jwt"
	"github.com/maximfedotov74/diploma-backend/internal/shared/password"
)

type userMailService interface {
	SendChangePasswordEmail(to string, subject string, code string) error
}

type userSessionServie interface {
	Create(ctx context.Context, dto model.CreateSessionDto) fall.Error
	Sign(claims jwt.UserClaims) (*jwt.Tokens, fall.Error)
	GetUserSessions(ctx context.Context, userId int, token string) (*model.UserSessionsResponse, fall.Error)

	RemoveSession(ctx context.Context, userId int, sessionId int) fall.Error
	RemoveExceptCurrentSession(ctx context.Context, userId int, sessionId int) fall.Error
	RemoveAllSessions(ctx context.Context, userId int) fall.Error
}

type userRepository interface {
	Create(ctx context.Context, dto model.CreateUserDto) (*model.CreatedUserResponse, fall.Error)
	FindByEmail(ctx context.Context, email string) (*model.User, fall.Error)
	FindById(ctx context.Context, id int) (*model.User, fall.Error)
	Update(ctx context.Context, dto model.UpdateUserDto, id int) fall.Error
	CreateChangePasswordCode(ctx context.Context, userId int) (*string, fall.Error)
	FindChangePasswordCode(ctx context.Context, userId int, code string) (*model.ChangePasswordCode, fall.Error)
	RemoveChangePasswordCode(ctx context.Context, userId int, tx db.Transaction) error
	ChangePassword(ctx context.Context, userId int, newPassword string) fall.Error
}

type UserService struct {
	repo           userRepository
	sessionService userSessionServie
	mailService    userMailService
}

func NewUserService(repo userRepository, sessionService userSessionServie,
	mailService userMailService) *UserService {
	return &UserService{repo: repo, sessionService: sessionService, mailService: mailService}
}

func (s *UserService) RemoveAllSessions(ctx context.Context, userId int) fall.Error {
	return s.sessionService.RemoveAllSessions(ctx, userId)
}
func (s *UserService) RemoveSession(ctx context.Context, userId int, sessionId int) fall.Error {
	return s.sessionService.RemoveSession(ctx, userId, sessionId)
}

func (s *UserService) RemoveExceptCurrentSession(ctx context.Context, userId int, sessionId int) fall.Error {
	return s.sessionService.RemoveExceptCurrentSession(ctx, userId, sessionId)
}

func (s *UserService) Create(ctx context.Context, dto model.CreateUserDto) (*model.CreatedUserResponse, fall.Error) {
	return s.repo.Create(ctx, dto)
}

func (s *UserService) FindById(ctx context.Context, id int) (*model.User, fall.Error) {
	return s.repo.FindById(ctx, id)
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*model.User, fall.Error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *UserService) Update(ctx context.Context, dto model.UpdateUserDto, id int) fall.Error {
	return s.repo.Update(ctx, dto, id)
}

func (us *UserService) CreateChangePasswordCode(ctx context.Context, userId int) fall.Error {

	currentUser, appErr := us.FindById(ctx, userId)
	if appErr != nil {
		return appErr
	}

	code, err := us.repo.CreateChangePasswordCode(ctx, currentUser.Id)

	if err != nil {
		return err
	}

	go us.mailService.SendChangePasswordEmail(currentUser.Email, "Код для смены пароля", *code)

	return nil
}

func (us *UserService) ConfirmChangePassword(ctx context.Context, code string, userId int) fall.Error {
	_, ex := us.repo.FindChangePasswordCode(ctx, userId, code)

	if ex != nil {
		return ex
	}

	err := us.repo.RemoveChangePasswordCode(ctx, userId, nil)

	if err != nil {
		return fall.ServerError(err.Error())
	}

	return nil

}

func (us *UserService) ChangePassword(ctx context.Context, dto model.ChangePasswordDto, localSession model.LocalSession) (*jwt.Tokens, fall.Error) {

	user, appErr := us.FindById(ctx, localSession.UserId)

	if appErr != nil {
		return nil, appErr
	}

	oldMatch := password.ComparePasswords(user.PasswordHash, dto.OldPassword)

	if !oldMatch {
		return nil, fall.NewErr(msg.BadPassword, fall.STATUS_BAD_REQUEST)
	}

	newMatch := password.ComparePasswords(user.PasswordHash, dto.NewPassword)

	if newMatch {
		return nil, fall.NewErr(msg.BadNewPassword, fall.STATUS_BAD_REQUEST)

	}

	newHash, err := password.HashPassword(dto.NewPassword)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	ex := us.repo.ChangePassword(ctx, user.Id, newHash)

	if ex != nil {
		return nil, ex
	}

	tokens, ex := us.sessionService.Sign(jwt.UserClaims{UserId: user.Id, UserAgent: localSession.UserAgent})

	if ex != nil {
		return nil, ex
	}

	tokenDto := model.CreateSessionDto{UserId: user.Id, UserAgent: localSession.UserAgent, Token: tokens.RefreshToken}
	appErr = us.sessionService.Create(ctx, tokenDto)

	if appErr != nil {
		return nil, appErr
	}

	return tokens, nil
}

func (s *UserService) GetUserSessions(ctx context.Context, userId int, token string) (*model.UserSessionsResponse, fall.Error) {
	return s.sessionService.GetUserSessions(ctx, userId, token)
}
