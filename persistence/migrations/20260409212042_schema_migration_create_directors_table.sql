-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS directors (
    director_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS directors;
-- +goose StatementEnd
