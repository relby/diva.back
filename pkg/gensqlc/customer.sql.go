// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: customer.sql

package gensqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const insertCustomer = `-- name: InsertCustomer :one
INSERT INTO customers (full_name, phone_number, discount)
VALUES ($1, $2, $3)
RETURNING id
`

type InsertCustomerParams struct {
	FullName    string
	PhoneNumber string
	Discount    int16
}

func (q *Queries) InsertCustomer(ctx context.Context, arg InsertCustomerParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertCustomer, arg.FullName, arg.PhoneNumber, arg.Discount)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const selectCustomerById = `-- name: SelectCustomerById :one
SELECT id, full_name, phone_number, discount FROM customers
WHERE id = $1
LIMIT 1
`

func (q *Queries) SelectCustomerById(ctx context.Context, id int64) (*Customer, error) {
	row := q.db.QueryRow(ctx, selectCustomerById, id)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.PhoneNumber,
		&i.Discount,
	)
	return &i, err
}

const selectCustomerByIdForUpdate = `-- name: SelectCustomerByIdForUpdate :one
SELECT id, full_name, phone_number, discount FROM customers
WHERE id = $1
LIMIT 1
FOR UPDATE
`

func (q *Queries) SelectCustomerByIdForUpdate(ctx context.Context, id int64) (*Customer, error) {
	row := q.db.QueryRow(ctx, selectCustomerByIdForUpdate, id)
	var i Customer
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.PhoneNumber,
		&i.Discount,
	)
	return &i, err
}

const selectCustomers = `-- name: SelectCustomers :many
SELECT id, full_name, phone_number, discount FROM customers
WHERE (full_name ILIKE '%' || $1 || '%' OR $1 IS NULL) AND (phone_number ILIKE '%' || $2 || '%' OR $2 IS NULL)
`

type SelectCustomersParams struct {
	FullName    pgtype.Text
	PhoneNumber pgtype.Text
}

func (q *Queries) SelectCustomers(ctx context.Context, arg SelectCustomersParams) ([]*Customer, error) {
	rows, err := q.db.Query(ctx, selectCustomers, arg.FullName, arg.PhoneNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Customer
	for rows.Next() {
		var i Customer
		if err := rows.Scan(
			&i.ID,
			&i.FullName,
			&i.PhoneNumber,
			&i.Discount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const truncateCustomers = `-- name: TruncateCustomers :exec
TRUNCATE customers RESTART IDENTITY CASCADE
`

func (q *Queries) TruncateCustomers(ctx context.Context) error {
	_, err := q.db.Exec(ctx, truncateCustomers)
	return err
}

const upsertCustomer = `-- name: UpsertCustomer :one
INSERT INTO customers (id, full_name, phone_number, discount)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id)
DO UPDATE SET full_name = $2, phone_number = $3, discount = $4
RETURNING id
`

type UpsertCustomerParams struct {
	ID          int64
	FullName    string
	PhoneNumber string
	Discount    int16
}

func (q *Queries) UpsertCustomer(ctx context.Context, arg UpsertCustomerParams) (int64, error) {
	row := q.db.QueryRow(ctx, upsertCustomer,
		arg.ID,
		arg.FullName,
		arg.PhoneNumber,
		arg.Discount,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}