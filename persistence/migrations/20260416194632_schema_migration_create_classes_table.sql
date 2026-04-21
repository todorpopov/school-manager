-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS classes (
    class_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    grade_level INTEGER NOT NULL,
    class_name VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS classes;
-- +goose StatementEnd
