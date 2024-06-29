-- name: SelectEmployees :many
SELECT sqlc.embed(users), sqlc.embed(employees) FROM employees
INNER JOIN users ON employees.user_id = users.id;

-- name: SelectEmployeeByID :one
SELECT sqlc.embed(users), sqlc.embed(employees) FROM employees
INNER JOIN users ON employees.user_id = users.id
WHERE users.id = $1
LIMIT 1;

-- name: SelectEmployeeByAccessKey :one
SELECT sqlc.embed(users), sqlc.embed(employees) FROM employees
INNER JOIN users ON employees.user_id = users.id
WHERE employees.access_key = $1
LIMIT 1;

-- name: UpsertEmployee :exec
INSERT INTO employees (user_id, access_key, permissions)
VALUES ($1, $2, $3)
ON CONFLICT (user_id)
DO UPDATE SET access_key = $2, permissions = $3;