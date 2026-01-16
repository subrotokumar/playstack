package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewSQLStore(connPool *pgxpool.Pool) *SQLStore {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
