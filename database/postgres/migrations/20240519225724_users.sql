-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    full_name text NOT NULL,

    PRIMARY KEY (id)
);

CREATE TABLE admins (
    user_id uuid NOT NULL,
    login text NOT NULL UNIQUE,
    hashed_password text NOT NULL,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
            ON DELETE CASCADE
            ON UPDATE CASCADE,

    PRIMARY KEY (user_id)
);

CREATE TYPE employee_permission AS ENUM('CREATE', 'UPDATE', 'DELETE');

CREATE TABLE employees (
    user_id uuid NOT NULL,
    access_key text NOT NULL UNIQUE,
    permissions employee_permission[] NOT NULL DEFAULT ARRAY[]::employee_permission[],

    FOREIGN KEY (user_id)
        REFERENCES users(id)
            ON DELETE CASCADE
            ON UPDATE CASCADE,

    PRIMARY KEY (user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE employees;
DROP TYPE employee_permission;
DROP TABLE admins;
DROP TABLE users;
-- +goose StatementEnd
