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
    cognito_sub,
    email
) VALUES (
    $1, $2, $3
)
RETURNING id, cognito_sub, email, created_at
`

type CreateUserParams struct {
	ID         uuid.UUID   `json:"id"`
	CognitoSub string      `json:"cognito_sub"`
	Email      pgtype.Text `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.ID, arg.CognitoSub, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CognitoSub,
		&i.Email,
		&i.CreatedAt,
	)
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

const getUserByCognitoSub = `-- name: GetUserByCognitoSub :one
SELECT id, cognito_sub, email, created_at
FROM users
WHERE cognito_sub = $1
`

func (q *Queries) GetUserByCognitoSub(ctx context.Context, cognitoSub string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByCognitoSub, cognitoSub)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CognitoSub,
		&i.Email,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, cognito_sub, email, created_at
FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email pgtype.Text) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CognitoSub,
		&i.Email,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, cognito_sub, email, created_at
FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CognitoSub,
		&i.Email,
		&i.CreatedAt,
	)
	return i, err
}
