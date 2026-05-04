-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS terms (
    term_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS terms;
-- +goose StatementEnd
