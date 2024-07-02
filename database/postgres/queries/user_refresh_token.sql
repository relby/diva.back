-- name: SelectUserRefreshTokenById :one
SELECT * FROM user_refresh_tokens
WHERE id = $1
LIMIT 1;

-- name: UpsertUserRefreshToken :exec
INSERT INTO user_refresh_tokens (id, user_id, expires_at)
VALUES ($1, $2, $3)
ON CONFLICT (id)
DO UPDATE SET user_id = $2, expires_at = $3;

-- name: DeleteRefreshTokenById :exec
DELETE FROM user_refresh_tokens WHERE id = $1;