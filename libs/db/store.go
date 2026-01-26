package db

import (
	"context"

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

func NewPgxPool(conn string, minConn, maxConn int32) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, err
	}
	config.MaxConns = maxConn
	config.MinConns = minConn
	return pgxpool.NewWithConfig(context.Background(), config)
}

func (store *SQLStore) Close() {
	store.connPool.Close()
}

func (store *SQLStore) Ping(ctx context.Context) error {
	return store.connPool.Ping(ctx)
}

func (store *SQLStore) Acquire(ctx context.Context) (c *pgxpool.Conn, err error) {
	return store.connPool.Acquire(ctx)
}

func (store *SQLStore) Stat() *pgxpool.Stat {
	return store.connPool.Stat()
}
