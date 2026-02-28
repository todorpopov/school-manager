-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
    role_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    role_name VARCHAR(255) NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
