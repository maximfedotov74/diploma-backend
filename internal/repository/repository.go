package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maximfedotov74/fiber-psql/internal/model"
)

type User interface {
	GetAll() error
	Create(dto model.CreateUserDto) (*UserRepoResponse, error)
	GetUserById(id int) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	FindActivationLink(link string) (*int, error)
	ActivateUser(id *int) error
}

type Role interface {
	Create(dto model.CreateRoleDto) (*model.Role, error)
	AddRoleToUser(roleId int, userId int) error
	RemoveRoleFromUser(roleId int, userId int) error
	FindRoleByTitle(title string) (*model.Role, error)
}

type Token interface {
	FindToken() error
	RemoveToken(token string) error
	UpdateToken() error
	CreateToken(dto model.CreateToken) error
}

type Repositories struct {
	UserRepository  User
	RoleRepository  Role
	TokenRepository Token
}

func New(db *pgxpool.Pool) *Repositories {

	role := NewRoleRepository(db)
	user := NewUserRepository(db)
	token := NewTokenRepository(db)

	return &Repositories{
		UserRepository:  user,
		RoleRepository:  role,
		TokenRepository: token,
	}
}
