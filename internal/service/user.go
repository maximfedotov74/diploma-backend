package service

import (
	"github.com/maximfedotov74/fiber-psql/internal/constants"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/repository"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
	"github.com/maximfedotov74/fiber-psql/pkg/token"
)

type UserService struct {
	repo            repository.User
	tokenService    Token
	passwordService Password
	mailService     Mail
}

func NewUserService(repo repository.User, tokenService Token, mailService Mail, passwordService Password) *UserService {
	return &UserService{
		repo:            repo,
		tokenService:    tokenService,
		mailService:     mailService,
		passwordService: passwordService,
	}
}

func (us *UserService) GetAll() {}

func (us *UserService) Create(dto model.CreateUserDto) (*model.UserCreatedResponse, lib.Error) {

	response, err := us.repo.Create(&dto.Password, dto.Email, constants.CREDENTIALS)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	return response, nil

}

func (us *UserService) CreateYandex(email string) (*model.UserCreatedResponse, lib.Error) {

	response, err := us.repo.Create(nil, email, constants.YANDEX)

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

func (us *UserService) CreateChangePasswordCode(user model.User) lib.Error {

	code, err := us.repo.CreateChangePasswordCode(user.Id)

	if err != nil {
		return lib.NewErr(err.Error(), 500)
	}

	go us.mailService.SendChangePasswordEmail(user.Email, "Код для смены пароля", *code)

	return nil
}

func (us *UserService) ChangePassword(dto model.ChangePasswordDto, contextData *model.UserContextData) (*token.Tokens, lib.Error) {

	if contextData.User.PasswordHash == nil {
		return nil, lib.NewErr(messages.BAD_NEW_PASSWORD, 400)
	}

	oldMatch := us.passwordService.ComparePasswords(*contextData.User.PasswordHash, dto.OldPassword)

	if !oldMatch {
		return nil, lib.NewErr(messages.BAD_PASSWORD, 400)
	}

	newMatch := us.passwordService.ComparePasswords(*contextData.User.PasswordHash, dto.NewPassword)

	if newMatch {
		return nil, lib.NewErr(messages.BAD_NEW_PASSWORD, 400)

	}

	newHash, err := us.passwordService.HashPassword(dto.NewPassword)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	_, err = us.repo.FindChangePasswordCode(contextData.User.Id, dto.Code)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 404)
	}

	err = us.repo.ChangePassword(contextData.User.Id, newHash)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	tokens, err := us.tokenService.Sign(token.UserClaims{UserId: contextData.User.Id, UserAgent: contextData.UserAgent})

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	tokenDto := model.CreateToken{UserId: contextData.User.Id, UserAgent: contextData.UserAgent, Token: tokens.RefreshToken}
	appErr := us.tokenService.Create(tokenDto)

	if appErr != nil {
		return nil, appErr
	}

	return &tokens, nil
}
