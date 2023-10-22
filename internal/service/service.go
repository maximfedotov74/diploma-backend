package service

import (
	"github.com/maximfedotov74/fiber-psql/internal/cfg"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/internal/repository"
	"github.com/maximfedotov74/fiber-psql/pkg/lib"
	"github.com/maximfedotov74/fiber-psql/pkg/token"
)

type User interface {
	GetAll()
	Create(dto model.CreateUserDto) (int, lib.Error)
	Login(dto model.LoginDto) (*model.LoginResponse, lib.Error)
	GetUserById(id int) (*model.User, lib.Error)
}

type Role interface {
	Create(dto model.CreateRoleDto) (*model.Role, error)
	AddRoleToUser(title string, userId int) (bool, error)
	RemoveRoleFromUser(title string, userId int) (bool, error)
}

type Services struct {
	UserService  User
	TokenService token.Token
	RoleService  Role
}

type Deps struct {
	Repos  *repository.Repositories
	Config *cfg.Config
}

func New(deps Deps) *Services {

	tokenService := token.New(deps.Config)
	roleService := NewRoleService(deps.Repos.RoleRepository)
	userService := NewUserService(deps.Repos.UserRepository, tokenService)
	return &Services{
		UserService:  userService,
		TokenService: tokenService,
		RoleService:  roleService,
	}
}
