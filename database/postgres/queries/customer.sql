-- name: SelectCustomers :many
SELECT * FROM customers
WHERE (full_name ILIKE '%' || sqlc.narg(full_name) || '%' OR sqlc.narg(full_name) IS NULL) AND (phone_number ILIKE '%' || sqlc.narg(phone_number) || '%' OR sqlc.narg(phone_number) IS NULL);

-- name: SelectCustomerById :one
SELECT * FROM customers
WHERE id = $1
LIMIT 1;

-- name: SelectCustomerByIdForUpdate :one
SELECT * FROM customers
WHERE id = $1
LIMIT 1
FOR UPDATE;

-- name: UpsertCustomer :one
INSERT INTO customers (id, full_name, phone_number, discount)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id)
DO UPDATE SET full_name = $2, phone_number = $3, discount = $4
RETURNING id;

-- name: InsertCustomer :one
INSERT INTO customers (full_name, phone_number, discount)
VALUES ($1, $2, $3)
RETURNING id;

-- name: TruncateCustomers :exec
TRUNCATE customers RESTART IDENTITY CASCADE;