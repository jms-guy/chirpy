-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: ClearUsers :exec
DELETE FROM users;

-- name: GetUserFromEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserFromID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUserInfo :exec
UPDATE users
SET email = $1, hashed_password = $2, updated_at = NOW()
WHERE id = $3;