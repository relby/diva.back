-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_refresh_tokens (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    expires_at timestamptz NOT NULL,

    FOREIGN KEY (user_id)
        REFERENCES users(id)
            ON DELETE CASCADE
            ON UPDATE CASCADE,

    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_refresh_tokens;
-- +goose StatementEnd
