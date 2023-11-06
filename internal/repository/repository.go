package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maximfedotov74/fiber-psql/internal/model"
)

const (
	emailField  = "email"
	userIdField = "user_id"
)

const (
	userRole  = "USER"
	adminRole = "ADMIN"
)

type User interface {
	GetAll() error
	Create(password string, email string) (*model.UserCreatedResponse, error)
	GetUserById(id int) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	FindActivationLink(link string) (*int, error)
	ActivateUser(id *int) error
	ChangePassword(userId int, newPassword string) error
	CreateChangePasswordCode(userId int) (*string, error)
	FindChangePasswordCode(userId int, code string) (*model.ChangePasswordCode, error)
}

type Role interface {
	Create(dto model.CreateRoleDto) (*model.Role, error)
	AddRoleToUser(roleId int, userId int) error
	RemoveRoleFromUser(roleId int, userId int) error
	FindRoleByTitle(title string) (*model.Role, error)
}

type Token interface {
	FindByAgentAndToken(agent string, token string) (*model.Token, error)
	RemoveToken(token string) error
	CreateToken(dto model.CreateToken) error
}

type Category interface {
	CreateCategoryType(title string) error
	CreateCategory(title string, img *string, parentId *int) error
	FindTypeByTitle(title string) (*model.CategoryType, error)
	FindCategoryByTitle(title string) (*model.Category, error)
}

type Repositories struct {
	UserRepository     User
	RoleRepository     Role
	TokenRepository    Token
	CategoryRepository Category
}

func New(db *pgxpool.Pool) *Repositories {

	role := NewRoleRepository(db)
	user := NewUserRepository(db)
	token := NewTokenRepository(db)
	categoryRepository := NewCategoryRepository(db)

	return &Repositories{
		UserRepository:     user,
		RoleRepository:     role,
		TokenRepository:    token,
		CategoryRepository: categoryRepository,
	}
}
