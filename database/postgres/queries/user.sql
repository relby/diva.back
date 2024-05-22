-- name: UpsertUser :exec
INSERT INTO users (id, full_name)
VALUES ($1, $2)
ON CONFLICT (id)
DO UPDATE SET full_name = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;