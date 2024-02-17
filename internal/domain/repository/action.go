package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/maximfedotov74/diploma-backend/internal/domain/msg"
	"github.com/maximfedotov74/diploma-backend/internal/shared/db"
	"github.com/maximfedotov74/diploma-backend/internal/shared/fall"
)

type ActionRepository struct {
	db db.PostgresClient
}

func NewActionRepository(db db.PostgresClient) *ActionRepository {
	return &ActionRepository{db: db}
}

func (r *ActionRepository) Create(ctx context.Context, dto model.CreateActionDto) fall.Error {
	q := "INSERT INTO action (end_date,title,img_path,description) VALUES ($1,$2,$3,$4);"

	_, err := r.db.Exec(ctx, q, dto.EndDate, dto.Title, dto.ImgPath, dto.Description)
	if err != nil {
		return fall.ServerError(fmt.Sprintf("Ошибка при создании акции: %s", err.Error()))
	}

	return nil
}

func (r *ActionRepository) AddModel(ctx context.Context, actionId string, modelId int) fall.Error {
	q := "INSERT INTO action_model (action_id,product_model_id) VALUES ($1,$2);"

	_, err := r.db.Exec(ctx, q, actionId, modelId)
	if err != nil {
		return fall.ServerError(fmt.Sprintf("Ошибка при добавлении модели в акцию: %s", err.Error()))
	}
	return nil
}

func (r *ActionRepository) FindById(ctx context.Context, id string) (*model.Action, fall.Error) {
	q := "SELECT action_id, created_at, updated_at, end_date, title, is_activated, img_path, description FROM action WHERE action_id = $1;"

	row := r.db.QueryRow(ctx, q, id)

	action := model.Action{}

	err := row.Scan(&action.Id, &action.CreatedAt, &action.UpdatedAt, &action.EndDate, &action.Title, &action.IsActivated,
		&action.ImgPath, &action.Description)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fall.NewErr(msg.ActionNotFound, fall.STATUS_NOT_FOUND)
		}
		return nil, fall.ServerError(err.Error())
	}

	return &action, nil
}

func (r *ActionRepository) GetAll(ctx context.Context) ([]model.Action, fall.Error) {
	q := "SELECT action_id, created_at, updated_at, end_date, title, is_activated, img_path, description FROM action ORDER BY created_at;"

	rows, err := r.db.Query(ctx, q)

	if err != nil {
		return nil, fall.ServerError(err.Error())
	}
	defer rows.Close()

	var actions []model.Action

	for rows.Next() {
		action := model.Action{}

		err := rows.Scan(&action.Id, &action.CreatedAt, &action.UpdatedAt, &action.EndDate, &action.Title, &action.IsActivated,
			&action.ImgPath, &action.Description)

		if err != nil {
			return nil, fall.ServerError(err.Error())
		}
		actions = append(actions, action)
	}

	if err := rows.Err(); err != nil {
		return nil, fall.ServerError(err.Error())
	}

	return actions, nil
}
