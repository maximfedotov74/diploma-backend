package service

import (
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/repository"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/mail"
	"github.com/maximfedotov74/fiber-psql/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo         repository.User
	tokenService token.Token
	mailService  mail.Mail
}

func NewUserService(repo repository.User, tokenService token.Token, mailService mail.Mail) *UserService {
	return &UserService{
		repo:         repo,
		tokenService: tokenService,
		mailService:  mailService,
	}
}

func (us *UserService) GetAll() {}
func (us *UserService) Create(dto model.CreateUserDto) (*int, lib.Error) {

	hash, err := us.hashPassword(dto.Password)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	dto.Password = hash

	response, err := us.repo.Create(dto)

	if err != nil {
		return &response.Id, lib.NewErr(err.Error(), 500)
	}

	go us.mailService.SendActivationEmail(dto.Email, "Активация аккаутна", response.ActivationAccountLink)

	// if err != nil {
	// 	return nil, lib.NewErr("Ошибка при отправке письма подтверждения на эл. почту!", 500)
	// }

	return &response.Id, nil
}

func (us *UserService) Activate() {

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
		return nil, lib.NewErr("Неверный логин и пароль!", 404)
	}

	if isPasswordCorrect := us.comparePasswords(user.PasswordHash, dto.Password); !isPasswordCorrect {
		return nil, lib.NewErr("Неверный логин и пароль!", 404)
	}

	tokens, err := us.tokenService.Sign(user.Id)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	response := model.LoginResponse{Id: user.Id, Tokens: tokens, Roles: user.Roles}

	return &response, nil

}
