package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Transaction struct {
	Executer pgx.Tx
	Ctx      context.Context
}
