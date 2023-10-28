package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maximfedotov74/fiber-psql/internal/model"
	"github.com/maximfedotov74/fiber-psql/pkg/messages"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetAll() error {
	return nil
}

func (ur *UserRepository) Create(dto model.CreateUserDto) (*UserRepoResponse, error) {

	txCtx := context.Background()

	tx, err := ur.db.Begin(txCtx)
	defer func() {
		if err != nil {
			tx.Rollback(txCtx)
		} else {
			tx.Commit(txCtx)
		}
	}()

	if err != nil {
		return nil, nil
	}

	query := "INSERT INTO public.user (email, password_hash) VALUES ($1, $2) RETURNING user_id;"

	row := tx.QueryRow(context.Background(), query, dto.Email, dto.Password)
	var id int

	err = row.Scan(&id)
	if err != nil {
		return nil, err
	}

	role := model.Role{}

	rowrole := tx.QueryRow(txCtx, findRoleByTitle, userRole)
	err = rowrole.Scan(&role.Id, &role.Title)

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(txCtx, addRoleToUser, id, role.Id)

	if err != nil {
		return nil, err
	}

	query = "INSERT INTO public.user_settings (auth_provider, user_id, activation_account_link) VALUES ($1, $2, uuid_generate_v4()) RETURNING activation_account_link;"
	row = tx.QueryRow(txCtx, query, "credentials", id)
	var link string
	err = row.Scan(&link)
	if err != nil {
		return nil, err
	}

	return &UserRepoResponse{Id: id, ActivationAccountLink: link}, nil
}

func (ur *UserRepository) findByIdOrEmail(field string, value any) (*model.User, error) {
	query := fmt.Sprintf(`SELECT public.user.user_id, public.user.email, public.user.password_hash,
	role.title, role.role_id, public.user_settings.is_activated
	FROM public.user
	LEFT JOIN user_role ON public.user.user_id = user_role.user_id
	LEFT JOIN public.role ON public.role.role_id = user_role.role_id
	LEFT JOIN public.user_settings ON public.user.user_id = public.user_settings.user_id
	WHERE public.user.%s = $1;`, field)

	rows, err := ur.db.Query(context.Background(), query, value)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user model.User
	processedRows := 0
	for rows.Next() {
		role := model.Role{}
		err := rows.Scan(&user.Id, &user.Email, &user.PasswordHash, &role.Title, &role.Id, &user.IsActivated)
		if err != nil {
			return nil, err
		}
		user.Roles = append(user.Roles, role)
		processedRows++
	}

	if rows.Err() != nil {
		return nil, err
	}
	if processedRows == 0 {
		return nil, nil
	}

	return &user, nil

}

func (ur *UserRepository) GetUserById(id int) (*model.User, error) {
	user, err := ur.findByIdOrEmail(userIdField, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	user, err := ur.findByIdOrEmail(emailField, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) FindActivationLink(link string) (*int, error) {
	query := `SELECT "public"."user_settings".user_id FROM "public"."user_settings"
	WHERE "public".user_settings.activation_account_link = $1;`
	row := ur.db.QueryRow(context.Background(), query, link)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (ur *UserRepository) ActivateUser(id *int) error {
	query := `UPDATE "public".user_settings
	SET activation_account_link = NULL,
	is_activated = TRUE
	WHERE "public".user_settings.user_id = $1;`

	_, err := ur.db.Exec(context.Background(), query, id)
	if err != nil {
		return errors.New(messages.ACTIVATION_ERROR)
	}

	return nil
}
