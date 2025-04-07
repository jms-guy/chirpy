-- name: CreateToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;

-- name: GetToken :one
SELECT * FROM refresh_tokens
WHERE token = $1
AND revoked_at IS NULL;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = $1, updated_at = NOW()
WHERE token = $2;

-- name: GetUserFromToken :one
SELECT users.* FROM users
INNER JOIN refresh_tokens
ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1;