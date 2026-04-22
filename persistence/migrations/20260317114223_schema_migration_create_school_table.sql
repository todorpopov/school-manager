-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS schools (
    school_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    school_name VARCHAR(255) NOT NULL UNIQUE,
    school_address VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS schools;
-- +goose StatementEnd
