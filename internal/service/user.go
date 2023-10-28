package service

import (
	"fmt"

	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/repository"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/mail"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
	"github.com/maximfedotov74/fiber-psql/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo         repository.User
	tokenService Token
	mailService  mail.Mail
}

func NewUserService(repo repository.User, tokenService Token, mailService mail.Mail) *UserService {
	return &UserService{
		repo:         repo,
		tokenService: tokenService,
		mailService:  mailService,
	}
}

func (us *UserService) GetAll() {}
func (us *UserService) Create(dto model.CreateUserDto) (*int, lib.Error) {

	user, err := us.repo.GetUserByEmail(dto.Email)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	if user != nil {
		return nil, lib.NewErr(messages.USER_EXISTS, 400)
	}

	hash, err := us.hashPassword(dto.Password)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}
	dto.Password = hash

	response, err := us.repo.Create(dto)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	link := fmt.Sprintf("/api/user/activate/%s", response.ActivationAccountLink)

	go us.mailService.SendActivationEmail(dto.Email, "Активация аккаутна", link)

	return &response.Id, nil
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

	return user, nil
}

func (us *UserService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (us *UserService) comparePasswords(hashed string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
	if err != nil {
		return false
	}
	return true
}

func (us *UserService) Login(dto model.LoginDto) (*model.LoginResponse, lib.Error) {
	user, err := us.repo.GetUserByEmail(dto.Email)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	if user == nil {
		return nil, lib.NewErr(messages.USER_NOT_FOUND, 404)
	}

	if isPasswordCorrect := us.comparePasswords(user.PasswordHash, dto.Password); !isPasswordCorrect {
		return nil, lib.NewErr(messages.INVALID_CREDENTIALS, 404)
	}

	userAgent := "user-agent"

	tokens, err := us.tokenService.Sign(token.UserClaims{UserId: user.Id, UserAgent: userAgent})

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	tokenDto := model.CreateToken{UserId: user.Id, UserAgent: userAgent, Token: tokens.RefreshToken}
	appErr := us.tokenService.Create(tokenDto)
	if appErr != nil {
		return nil, appErr
	}

	response := model.LoginResponse{Id: user.Id, Tokens: tokens, Roles: user.Roles}

	return &response, nil

}

func (us *UserService) GetLk(id int) (*model.User, lib.Error) {
	user, err := us.repo.GetUserById(id)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	if user == nil {
		return nil, lib.NewErr(messages.USER_NOT_FOUND, 404)
	}

	return user, nil
}
