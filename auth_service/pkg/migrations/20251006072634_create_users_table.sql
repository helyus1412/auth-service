-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    email varchar(25) NOT NULL,
    password varchar NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    created_by VARCHAR(255),
    updated_at TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(255)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users IF EXISTS;
-- +goose StatementEnd
