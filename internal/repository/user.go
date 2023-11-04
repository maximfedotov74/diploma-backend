package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
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

func (ur *UserRepository) Create(dto model.CreateUserDto) (*model.UserCreatedResponse, error) {

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
		return nil, err
	}

	query := "INSERT INTO public.user (email, password_hash) VALUES ($1, $2) RETURNING user_id;"

	row := tx.QueryRow(txCtx, query, dto.Email, dto.Password)
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

	return &model.UserCreatedResponse{Id: id, ActivationAccountLink: link}, nil
}

func (ur *UserRepository) findByIdOrEmail(field string, value any) (*model.User, error) {

	ctx := context.Background()

	query := fmt.Sprintf(`SELECT public.user.user_id, public.user.email, public.user.password_hash,
	role.title, role.role_id, public.user_settings.is_activated
	FROM public.user
	LEFT JOIN user_role ON public.user.user_id = user_role.user_id
	LEFT JOIN public.role ON public.role.role_id = user_role.role_id
	LEFT JOIN public.user_settings ON public.user.user_id = public.user_settings.user_id
	WHERE public.user.%s = $1;`, field)

	rows, err := ur.db.Query(ctx, query, value)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user := model.User{}
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
	ctx := context.Background()
	query := `SELECT "public"."user_settings".user_id FROM "public"."user_settings"
	WHERE "public".user_settings.activation_account_link = $1;`
	row := ur.db.QueryRow(ctx, query, link)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (ur *UserRepository) ActivateUser(id *int) error {
	ctx := context.Background()
	query := `UPDATE "public".user_settings
	SET activation_account_link = NULL,
	is_activated = TRUE
	WHERE "public".user_settings.user_id = $1;`

	_, err := ur.db.Exec(ctx, query, id)
	if err != nil {
		return errors.New(messages.ACTIVATION_ERROR)
	}

	return nil
}

func (ur *UserRepository) ChangePassword(userId int, newPassword string) error {

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
		return err
	}

	query := `UPDATE public.user SET password_hash = $1,
	updated_at = CURRENT_TIMESTAMP WHERE public.user.user_id = $2;`

	_, err = tx.Exec(txCtx, query, newPassword, userId)
	if err != nil {
		return errors.New(messages.UPDATE_PASSWORD_ERROR)
	}

	err = RemoveChangePasswordCode(userId, tx, txCtx)
	if err != nil {
		return err
	}

	return nil
}

func RemoveChangePasswordCode(userId int, tx pgx.Tx, ctx context.Context) error {
	query := "DELETE FROM change_password_code WHERE user_id = $1"

	_, err := tx.Exec(ctx, query, userId)

	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) FindChangePasswordCode(userId int, code string) (*model.ChangePasswordCode, error) {
	ctx := context.Background()

	log.Println(userId)
	log.Println(code)
	query := `SELECT change_password_code_id, code, user_id FROM
	change_password_code WHERE user_id = $1 AND code = $2 AND end_time > CURRENT_TIMESTAMP;`

	row := ur.db.QueryRow(ctx, query, userId, code)

	codeModel := model.ChangePasswordCode{}

	err := row.Scan(&codeModel.ChangePasswordCodeId, &codeModel.Code, &codeModel.UserId)

	if err != nil {
		log.Println(err.Error())
		return nil, errors.New(messages.CHANGE_PASSWORD_CODE_NOT_FOUND)
	}

	return &codeModel, nil

}

func (ur *UserRepository) CreateChangePasswordCode(userId int) (*string, error) {

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
		return nil, err
	}

	err = RemoveChangePasswordCode(userId, tx, txCtx)

	if err != nil {
		return nil, err
	}

	query := "INSERT INTO change_password_code (user_id) VALUES ($1) RETURNING code;"

	row := tx.QueryRow(txCtx, query, userId)

	var code string

	err = row.Scan(&code)

	if err != nil {
		return nil, err
	}

	return &code, nil
}
