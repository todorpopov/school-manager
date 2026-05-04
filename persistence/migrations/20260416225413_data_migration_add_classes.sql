-- +goose Up
-- +goose StatementBegin
-- GPAE (school_id = 1) - Grades 1-4
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 1, 'A');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 1, 'B');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 1, 'C');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 1, 'D');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 1, 'F');

INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 2, 'A');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 2, 'B');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 2, 'C');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 2, 'D');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 2, 'F');

INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 3, 'A');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 3, 'B');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 3, 'C');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 3, 'D');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 3, 'F');

INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 4, 'A');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 4, 'B');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 4, 'C');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 4, 'D');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (1, 4, 'F');

-- SMG (school_id = 2) - Grades 5-8
INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 5, 'A');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 5, 'B');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 5, 'C');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 5, 'D');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 5, 'F');

INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 6, 'A');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 6, 'B');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 6, 'C');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 6, 'D');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 6, 'F');

INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 7, 'A');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 7, 'B');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 7, 'C');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 7, 'D');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (2, 7, 'F');

-- GPNE (school_id = 3) - Grades 8-12
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 8, 'A');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 8, 'B');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 8, 'C');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 8, 'D');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 8, 'F');

INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 9, 'A');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 9, 'B');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 9, 'C');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 9, 'D');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 9, 'F');

INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 10, 'A');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 10, 'B');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 10, 'C');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 10, 'D');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 10, 'F');

INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 11, 'A');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 11, 'B');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 11, 'C');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 11, 'D');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 11, 'F');

INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 12, 'A');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 12, 'B');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 12, 'C');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 12, 'D');
INSERT INTO classes (school_id, grade_level, class_name) VALUES (3, 12, 'F');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM classes;
-- +goose StatementEnd
