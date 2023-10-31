package service

import (
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/repository"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/mail"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
)

type UserService struct {
	repo            repository.User
	tokenService    Token
	passwordService Password
	mailService     mail.Mail
}

func NewUserService(repo repository.User, tokenService Token, mailService mail.Mail, passwordService Password) *UserService {
	return &UserService{
		repo:            repo,
		tokenService:    tokenService,
		mailService:     mailService,
		passwordService: passwordService,
	}
}

func (us *UserService) GetAll() {}

func (us *UserService) Create(dto model.CreateUserDto) (*model.UserCreatedResponse, lib.Error) {

	response, err := us.repo.Create(dto)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	return response, nil

}

func (us *UserService) Activate(activationLink string) lib.Error {
	id, err := us.repo.FindActivationLink(activationLink)

	if err != nil {
		return lib.NewErr(err.Error(), 404)
	}

	err = us.repo.ActivateUser(id)
	if err != nil {
		return lib.NewErr(err.Error(), 500)
	}

	return nil
}

func (us *UserService) GetUserById(id int) (*model.User, lib.Error) {
	user, err := us.repo.GetUserById(id)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	if user == nil {
		return nil, lib.NewErr(messages.USER_NOT_FOUND, 404)
	}

	return user, nil
}

func (us *UserService) GetUserByEmail(email string) (*model.User, lib.Error) {
	user, err := us.repo.GetUserByEmail(email)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	if user == nil {
		return nil, lib.NewErr(messages.USER_NOT_FOUND, 404)
	}

	return user, nil
}

func (us *UserService) ChangePassword(dto model.ChangePasswordDto, userId int, userAgent string) lib.Error {

	user, appErr := us.GetUserById(userId)

	if appErr != nil {
		return appErr
	}

	oldMatch := us.passwordService.ComparePasswords(user.PasswordHash, dto.OldPassword)

	if !oldMatch {
		return lib.NewErr(messages.BAD_PASSWORD, 400)
	}

	newMatch := us.passwordService.ComparePasswords(user.PasswordHash, dto.NewPassword)

	if newMatch {
		return lib.NewErr(messages.BAD_NEW_PASSWORD, 400)

	}

	newHash, err := us.passwordService.HashPassword(dto.NewPassword)

	if err != nil {
		return lib.NewErr(err.Error(), 500)
	}

	err = us.repo.ChangePassword(user.Id, newHash)

	if err != nil {
		return lib.NewErr(err.Error(), 500)
	}

	//update tokens

	return nil
}
