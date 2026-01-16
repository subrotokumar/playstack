-- name: CreateUser :one
INSERT INTO users (
    id,
    cognito_sub,
    email
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByCognitoSub :one
SELECT *
FROM users
WHERE cognito_sub = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;