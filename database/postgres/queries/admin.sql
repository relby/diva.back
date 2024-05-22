-- name: UpsertAdmin :exec
INSERT INTO admins (user_id, login, hashed_password)
VALUES ($1, $2, $3)
ON CONFLICT (user_id)
DO UPDATE SET login = $2, hashed_password = $3;