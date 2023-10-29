package service

import (
	"github.com/maximfedotov74/fiber-psql/internal/cfg"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/repository"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/mail"
	"github.com/maximfedotov74/fiber-psql/pkg/token"
)

type User interface {
	GetAll()
	Create(dto model.CreateUserDto) (*int, lib.Error)
	Login(dto model.LoginDto, userAgent string) (*model.LoginResponse, lib.Error)
	GetUserById(id int) (*model.User, lib.Error)
	Activate(activationLink string) lib.Error
	GetLk(id int) (*model.User, lib.Error)
	RefreshToken(refreshToken string, userAgent string) (*model.LoginResponse, lib.Error)
}

type Role interface {
	Create(dto model.CreateRoleDto) (*model.Role, lib.Error)
	AddRoleToUser(title string, userId int) lib.Error
	RemoveRoleFromUser(title string, userId int) lib.Error
}

type Token interface {
	FindToken(agent string, token string) (*model.Token, lib.Error)
	RemoveToken() error
	Sign(claims token.UserClaims) (token.Tokens, error)
	Parse(token string, tokenType token.TokenType) (*token.UserClaims, error)
	Refresh(refreshToken string) error
	Create(dto model.CreateToken) lib.Error
}

type Services struct {
	UserService  User
	RoleService  Role
	MailService  mail.Mail
	TokenService Token
}

type Deps struct {
	Repos  *repository.Repositories
	Config *cfg.Config
}

func New(deps Deps) *Services {
	mailService := mail.New(deps.Config.SmtpKey, deps.Config.SmtpMail, deps.Config.SmtpHost, deps.Config.SmtpPort, deps.Config.AppLink)
	tokenService := NewTokenService(deps.Config, deps.Repos.TokenRepository)
	roleService := NewRoleService(deps.Repos.RoleRepository)
	userService := NewUserService(deps.Repos.UserRepository, tokenService, mailService)
	return &Services{
		UserService:  userService,
		TokenService: tokenService,
		RoleService:  roleService,
	}
}
