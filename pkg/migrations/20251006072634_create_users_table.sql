-- +goose Up
-- +goose StatementBegin
CREATE TABLE users IF NOT EXISTS (
    id serial PRIMARY KEY,
    email varchar(25) NOT NULL,
    password varchar NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users IF EXISTS
-- +goose StatementEnd
