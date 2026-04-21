-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS grades (
    grade_id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL REFERENCES students(student_id) ON DELETE RESTRICT,
    curriculum_id INTEGER NOT NULL REFERENCES curricula(curriculum_id) ON DELETE RESTRICT,
    grade_value NUMERIC(3,2) NOT NULL CHECK (grade_value >= 2.00 AND grade_value <= 6.00),
    grade_date DATE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS grades;
-- +goose StatementEnd