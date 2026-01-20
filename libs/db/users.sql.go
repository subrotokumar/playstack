// source: users.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    id,
    email
) VALUES (
    $1, $2
)
RETURNING id, email, created_at
`

type CreateUserParams struct {
	ID    uuid.UUID   `json:"id"`
	Email pgtype.Text `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.ID, arg.Email)
	var i User
	err := row.Scan(&i.ID, &i.Email, &i.CreatedAt)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, created_at
FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email pgtype.Text) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(&i.ID, &i.Email, &i.CreatedAt)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, email, created_at
FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(&i.ID, &i.Email, &i.CreatedAt)
	return i, err
}
