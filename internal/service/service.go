package service

import (
	"github.com/maximfedotov74/fiber-psql/internal/cfg"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/repository"
	"github.com/maximfedotov74/fiber-psql/pkg/ip"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/mail"
	"github.com/maximfedotov74/fiber-psql/pkg/token"
)

type User interface {
	GetAll()
	Create(dto model.CreateUserDto) (*model.UserCreatedResponse, lib.Error)
	GetUserById(id int) (*model.User, lib.Error)
	GetUserByEmail(email string) (*model.User, lib.Error)
	Activate(activationLink string) lib.Error
	CreateChangePasswordCode(user model.User) lib.Error
	ChangePassword(dto model.ChangePasswordDto, contextData *model.UserContextData) (*token.Tokens, lib.Error)
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

type Auth interface {
	Registration(dto model.CreateUserDto) (*int, lib.Error)
	Login(dto model.LoginDto, userAgent string) (*model.LoginResponse, lib.Error)
	Refresh(refreshToken string, userAgent string) (*model.LoginResponse, lib.Error)
}

type Password interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashed string, pass string) bool
}

type Mail interface {
	SendActivationEmail(to string, subject string, link string) error
	SendChangePasswordEmail(to string, subject string, code string) error
}

type Ip interface {
	GetGeolocation(ip string) (*ip.IpLocationResponse, error)
}

type Category interface {
	CreateCategoryType(dto model.CreateCategoryTypeDto) lib.Error
	FindTypeByTitle(title string) (*model.CategoryType, lib.Error)
	FindCategoryByTitle(title string) (*model.Category, lib.Error)
	CreateCategory(dto model.CreateCategoryDto) lib.Error
}

type Services struct {
	UserService     User
	RoleService     Role
	MailService     Mail
	TokenService    Token
	AuthService     Auth
	CategoryService Category
}

type Deps struct {
	Repos  *repository.Repositories
	Config *cfg.Config
}

func New(deps Deps) *Services {
	passwordService := NewPasswordService()
	mailService := mail.New(deps.Config.SmtpKey, deps.Config.SmtpMail, deps.Config.SmtpHost, deps.Config.SmtpPort, deps.Config.AppLink)
	tokenService := NewTokenService(deps.Config, deps.Repos.TokenRepository)
	roleService := NewRoleService(deps.Repos.RoleRepository)
	userService := NewUserService(deps.Repos.UserRepository, tokenService, mailService, passwordService)
	authService := NewAuthService(userService, tokenService, passwordService, mailService)
	categoryService := NewCategoryService(deps.Repos.CategoryRepository)

	return &Services{
		UserService:     userService,
		TokenService:    tokenService,
		RoleService:     roleService,
		AuthService:     authService,
		CategoryService: categoryService,
	}
}
