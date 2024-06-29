-- name: UpsertAdmin :exec
INSERT INTO admins (user_id, login, hashed_password)
VALUES ($1, $2, $3)
ON CONFLICT (user_id)
DO UPDATE SET login = $2, hashed_password = $3;

-- name: SelectAdminById :one
SELECT sqlc.embed(users), sqlc.embed(admins) FROM admins
INNER JOIN users ON admins.user_id = users.id
WHERE users.id = $1
LIMIT 1;

-- name: SelectAdminByLogin :one
SELECT sqlc.embed(users), sqlc.embed(admins) FROM admins
INNER JOIN users ON admins.user_id = users.id
WHERE admins.login = $1
LIMIT 1;