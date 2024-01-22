package db

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbInstance *pgxpool.Pool
var pgOnce sync.Once

func NewPostgresConnection(connString string) *pgxpool.Pool {
	pgOnce.Do(func() {
		db, err := pgxpool.New(context.Background(), connString)
		if err != nil {
			log.Fatalf("Cannot connect to db, error: %s", err.Error())
		}
		dbInstance = db
	})

	return dbInstance
}
