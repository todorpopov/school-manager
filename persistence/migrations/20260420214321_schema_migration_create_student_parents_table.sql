-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS student_parents (
    parent_id INTEGER NOT NULL REFERENCES parents(parent_id) ON DELETE CASCADE,
    student_id INTEGER NOT NULL REFERENCES students(student_id) ON DELETE CASCADE,
    PRIMARY KEY (parent_id, student_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_parents;
-- +goose StatementEnd
