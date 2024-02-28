package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

const USER_ROLE = "USER"

type userRoleRepository interface {
	FindRoleByTitle(ctx context.Context, title string) (*model.Role, fall.Error)
	AddRoleToUser(ctx context.Context, roleId int, userId int, tx db.Transaction) fall.Error
}

type UserRepository struct {
	db       db.PostgresClient
	roleRepo userRoleRepository
}

func NewUserRepository(db db.PostgresClient, roleRepo userRoleRepository) *UserRepository {
	return &UserRepository{db: db, roleRepo: roleRepo}
}

func (r *UserRepository) Create(ctx context.Context, dto model.CreateUserDto) (*model.CreatedUserResponse, fall.Error) {

	tx, err := r.db.Begin(ctx)

	var ex fall.Error = nil

	if err != nil {
		ex = fall.ServerError(err.Error())
		return nil, ex
	}

	defer func() {
		if ex != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	q := "INSERT INTO public.user (email, password_hash) VALUES ($1,$2) RETURNING user_id, email;"

	row := tx.QueryRow(ctx, q, dto.Email, dto.Password)

	var email string
	var id int
	err = row.Scan(&id, &email)

	if err != nil {
		ex = fall.ServerError(err.Error())
		return nil, ex
	}

	role, ex := r.roleRepo.FindRoleByTitle(ctx, USER_ROLE)

	if ex != nil {
		return nil, ex
	}

	ex = r.roleRepo.AddRoleToUser(ctx, role.Id, id, tx)

	if ex != nil {
		return nil, ex
	}
	q = "INSERT INTO user_activation (user_id, activation_account_link) VALUES ($1, uuid_generate_v4()) RETURNING activation_account_link;"

	row = tx.QueryRow(ctx, q, id)

	var link string

	err = row.Scan(&link)

	if err != nil {
		ex = fall.ServerError(err.Error())
		return nil, ex
	}

	return &model.CreatedUserResponse{Id: id, Email: email, Link: link}, nil

}

func (r *UserRepository) Update(ctx context.Context, dto model.UpdateUserDto, id int) fall.Error {

	var queries []string

	if dto.AvatarPath != nil {
		queries = append(queries, fmt.Sprintf("avatar_path = '%s'", *dto.AvatarPath))
	}

	if dto.Patronymic != nil {
		queries = append(queries, fmt.Sprintf("patronymic = '%s'", *dto.Patronymic))
	}

	if dto.LastName != nil {
		queries = append(queries, fmt.Sprintf("last_name = '%s'", *dto.LastName))
	}

	if dto.FirstName != nil {
		queries = append(queries, fmt.Sprintf("first_name = '%s'", *dto.FirstName))
	}

	if dto.Gender != nil {
		queries = append(queries, fmt.Sprintf("gender = '%s'", *dto.Gender))
	}

	if len(queries) > 0 {
		q := "UPDATE public.user SET " + strings.Join(queries, ",") + " WHERE user_id = $1;"
		_, err := r.db.Exec(ctx, q, id)
		if err != nil {
			return fall.ServerError(fmt.Sprintf("%s, details: \n %s", msg.BrandUpdateError, err.Error()))
		}
		return nil
	}

	return nil
}

func (ur *UserRepository) findByIdOrEmail(ctx context.Context, field string, value any) (*model.User, fall.Error) {

	query := fmt.Sprintf(`
	SELECT public.user.user_id, public.user.email, public.user.password_hash,
	public.user.patronymic, public.user.first_name,	public.user.last_name, 
	role.title, role.role_id, public.user.is_activated, public.user.gender, public.user.avatar_path
	FROM public.user
	LEFT JOIN user_role ON public.user.user_id = user_role.user_id
	LEFT JOIN public.role ON public.role.role_id = user_role.role_id
	WHERE public.user.%s = $1;`, field)

	rows, err := ur.db.Query(ctx, query, value)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}

	defer rows.Close()

	user := model.User{}
	founded := false
	for rows.Next() {
		role := model.UserRole{}
		err := rows.Scan(&user.Id, &user.Email, &user.PasswordHash, &user.Patronymic, &user.FirstName, &user.LastName,
			&role.Title, &role.Id, &user.IsActivated, &user.Gender, &user.AvatarPath)
		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		user.Roles = append(user.Roles, role)
		if !founded {
			founded = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(rows.Err().Error())
	}

	if !founded {
		return nil, fall.NewErr(msg.UserNotFound, fall.STATUS_NOT_FOUND)
	}

	return &user, nil

}

func (r *UserRepository) FindById(ctx context.Context, id int) (*model.User, fall.Error) {
	return r.findByIdOrEmail(ctx, "user_id", id)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, fall.Error) {
	return r.findByIdOrEmail(ctx, "email", email)
}
