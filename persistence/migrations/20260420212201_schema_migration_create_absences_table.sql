-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS absences (
    absence_id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL REFERENCES students(student_id) ON DELETE RESTRICT,
    curriculum_id INTEGER NOT NULL REFERENCES curricula(curriculum_id) ON DELETE RESTRICT,
    absence_date DATE NOT NULL,
    is_excused BOOLEAN NOT NULL DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS absences;
-- +goose StatementEnd