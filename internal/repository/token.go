package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maximfedotov74/fiber-psql/internal/model"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{db: db}
}

func (tr *TokenRepository) FindByAgentAndToken(agent string, token string) (*model.Token, error) {

	ctx := context.Background()

	query := "SELECT token_id, user_id, user_agent, token FROM public.token WHERE user_agent = $1 AND token = $2;"

	tokenModel := model.Token{}

	row := tr.db.QueryRow(ctx, query, agent, token)

	err := row.Scan(&tokenModel.TokenId, &tokenModel.UserId, &tokenModel.UserAgent, &tokenModel.Token)

	if err != nil {
		return nil, err
	}

	return &tokenModel, nil

}

func (tr *TokenRepository) RemoveToken(token string) error {
	query := "DELETE FROM token WHERE token = $1;"
	_, err := tr.db.Exec(context.Background(), query, token)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TokenRepository) UpdateToken() error {
	return nil
}

func (tr *TokenRepository) CreateToken(dto model.CreateToken) error {

	ctx := context.Background()

	query := "SELECT token_id, user_id, user_agent, token FROM public.token WHERE user_id = $1 AND user_agent = $2;"

	tokenModel := model.Token{}

	row := tr.db.QueryRow(ctx, query, dto.UserId, dto.UserAgent)

	err := row.Scan(&tokenModel.TokenId, &tokenModel.UserId, &tokenModel.UserAgent, &tokenModel.Token)

	if err != nil {
		if err == pgx.ErrNoRows {

			query = "INSERT INTO token (token, user_agent, user_id) VALUES ($1, $2, $3) RETURNING token_id;"
			_, err = tr.db.Exec(ctx, query, dto.Token, dto.UserAgent, dto.UserId)

			if err != nil {
				return err
			}

			return nil
		} else {
			return err
		}
	}

	query = "UPDATE public.token SET token = $1, updated_at = CURRENT_TIMESTAMP;"

	_, err = tr.db.Exec(ctx, query, dto.Token)

	if err != nil {
		return err
	}

	return nil

}
