-- +goose Up
-- +goose StatementBegin
CREATE TABLE customers (
  id   bigserial NOT NULL,
  full_name text NOT NULL,
  phone_number text NOT NULL,
  discount smallint NOT NULL,

  PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE customers;
-- +goose StatementEnd
