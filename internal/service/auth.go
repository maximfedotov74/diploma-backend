package service

import (
	"fmt"

	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
	"github.com/maximfedotov74/fiber-psql/pkg/token"
)

type AuthService struct {
	userService     User
	tokenService    Token
	passwordService Password
	mailService     Mail
}

func NewAuthService(userService User, tokenService Token, passwordService Password, mailService Mail) *AuthService {
	return &AuthService{
		userService:     userService,
		tokenService:    tokenService,
		passwordService: passwordService,
		mailService:     mailService,
	}
}

func (as *AuthService) Login(dto model.LoginDto, userAgent string) (*model.LoginResponse, lib.Error) {
	user, appErr := as.userService.GetUserByEmail(dto.Email)

	if appErr != nil {
		return nil, appErr
	}

	if isPasswordCorrect := as.passwordService.ComparePasswords(user.PasswordHash, dto.Password); !isPasswordCorrect {
		return nil, lib.NewErr(messages.INVALID_CREDENTIALS, 404)
	}

	tokens, err := as.tokenService.Sign(token.UserClaims{UserId: user.Id, UserAgent: userAgent})

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	tokenDto := model.CreateToken{UserId: user.Id, UserAgent: userAgent, Token: tokens.RefreshToken}
	appErr = as.tokenService.Create(tokenDto)
	if appErr != nil {
		return nil, appErr
	}

	response := model.LoginResponse{Id: user.Id, Tokens: tokens, Roles: user.Roles}

	return &response, nil
}

func (as *AuthService) Registration(dto model.CreateUserDto) (*int, lib.Error) {
	user, _ := as.userService.GetUserByEmail(dto.Email)

	if user != nil {
		return nil, lib.NewErr(messages.USER_EXISTS, 400)
	}

	hash, err := as.passwordService.HashPassword(dto.Password)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}
	dto.Password = hash

	response, appErr := as.userService.Create(dto)

	if appErr != nil {
		return nil, appErr
	}

	link := fmt.Sprintf("/api/user/activate/%s", response.ActivationAccountLink)

	go as.mailService.SendActivationEmail(dto.Email, "Активация аккаутна", link)

	return &response.Id, nil

}

func (as *AuthService) Refresh(refreshToken string, userAgent string) (*model.LoginResponse, lib.Error) {

	_, err := as.tokenService.Parse(refreshToken, token.RefreshToken)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 401)
	}

	dbToken, appErr := as.tokenService.FindToken(userAgent, refreshToken)

	if appErr != nil {
		return nil, appErr
	}

	user, appErr := as.userService.GetUserById(dbToken.UserId)
	if appErr != nil {
		return nil, appErr
	}

	claims := token.UserClaims{UserId: user.Id, UserAgent: dbToken.UserAgent}

	tokens, err := as.tokenService.Sign(claims)

	if err != nil {
		return nil, lib.NewErr(err.Error(), 500)
	}

	tokenDto := model.CreateToken{UserId: user.Id, UserAgent: dbToken.UserAgent, Token: tokens.RefreshToken}
	appErr = as.tokenService.Create(tokenDto)

	if appErr != nil {
		return nil, appErr
	}

	response := model.LoginResponse{Id: user.Id, Tokens: tokens, Roles: user.Roles}

	return &response, nil
}
