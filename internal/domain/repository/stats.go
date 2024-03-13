package repository

import "github.com/maximfedotov74/diploma-backend/internal/shared/db"

type StatsRepository struct {
	db db.PostgresClient
}

func NewStatsRepository(db db.PostgresClient) *StatsRepository {
	return &StatsRepository{db: db}
}
