package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maximfedotov74/fiber-psql/internal/model"
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

func (ur *UserRepository) Create(dto model.CreateUserDto) (int, error) {
	var id int

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
		return id, nil
	}

	query := "INSERT INTO public.user (email, password_hash) VALUES ($1, $2) RETURNING user_id;"

	row := tx.QueryRow(context.Background(), query, dto.Email, dto.Password)

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	query = "INSERT INTO user_settings (auth_provider, user_id, activation_account_link) VALUES ($1, $2, uuid_generate_v4());"

	_, err = tx.Exec(txCtx, query, "credentials", id)
	if err != nil {
		return 0, err
	}

	role := model.Role{}

	rowrole := tx.QueryRow(txCtx, findRoleByTitle, userRole)
	err = rowrole.Scan(&role.Id, &role.Title)

	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(txCtx, addRoleToUser, id, role.Id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ur *UserRepository) findByIdOrEmail(field string, value any) (*model.User, error) {
	query := fmt.Sprintf(`SELECT public.user.user_id, public.user.email, public.user.password_hash, role.title, role.role_id FROM public.user
	LEFT JOIN user_role ON public.user.user_id = user_role.user_id
	LEFT JOIN public.role ON public.role.role_id = user_role.role_id
	WHERE public.user.%s = $1;`, field)

	rows, err := ur.db.Query(context.Background(), query, value)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := model.User{}

	for rows.Next() {
		role := model.Role{}
		err := rows.Scan(&user.Id, &user.Email, &user.PasswordHash, &role.Title, &role.Id)
		if err != nil {
			return nil, err
		}
		user.Roles = append(user.Roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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
