package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maximfedotov74/fiber-psql/internal/model"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{db: db}
}

func (tr *TokenRepository) FindToken() error {
	return nil
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
	log.Println(dto)
	query := "INSERT INTO token (token, user_agent, user_id) VALUES ($1, $2, $3) RETURNING token_id;"
	_, err := tr.db.Exec(context.Background(), query, dto.Token, dto.UserAgent, dto.UserId)

	if err != nil {
		return err
	}

	return nil
}
