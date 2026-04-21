-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS curricula (
    curriculum_id SERIAL PRIMARY KEY,
    class_id INTEGER NOT NULL REFERENCES classes(class_id) ON DELETE RESTRICT,
    subject_id INTEGER NOT NULL REFERENCES subjects(subject_id) ON DELETE RESTRICT,
    teacher_id INTEGER REFERENCES teachers(teacher_id) ON DELETE SET NULL,
    term_id INTEGER NOT NULL REFERENCES terms(term_id) ON DELETE RESTRICT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS curricula;
-- +goose StatementEnd
