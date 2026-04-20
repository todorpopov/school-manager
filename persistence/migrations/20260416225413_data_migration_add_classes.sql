-- +goose Up
-- +goose StatementBegin
INSERT INTO classes (grade_level, class_name) VALUES (1, 'A');
INSERT INTO classes (grade_level, class_name) VALUES (1, 'B');
INSERT INTO classes (grade_level, class_name) VALUES (1, 'C');
INSERT INTO classes (grade_level, class_name) VALUES (1, 'D');
INSERT INTO classes (grade_level, class_name) VALUES (1, 'F');

INSERT INTO classes (grade_level, class_name) VALUES (2, 'A');
INSERT INTO classes (grade_level, class_name) VALUES (2, 'B');
INSERT INTO classes (grade_level, class_name) VALUES (2, 'C');
INSERT INTO classes (grade_level, class_name) VALUES (2, 'D');
INSERT INTO classes (grade_level, class_name) VALUES (2, 'F');

INSERT INTO classes (grade_level, class_name) VALUES (3, 'A');
INSERT INTO classes (grade_level, class_name) VALUES (3, 'B');
INSERT INTO classes (grade_level, class_name) VALUES (3, 'C');
INSERT INTO classes (grade_level, class_name) VALUES (3, 'D');
INSERT INTO classes (grade_level, class_name) VALUES (3, 'F');

INSERT INTO classes (grade_level, class_name) VALUES (4, 'A');
INSERT INTO classes (grade_level, class_name) VALUES (4, 'B');
INSERT INTO classes (grade_level, class_name) VALUES (4, 'C');
INSERT INTO classes (grade_level, class_name) VALUES (4, 'D');
INSERT INTO classes (grade_level, class_name) VALUES (4, 'F');

INSERT INTO classes (grade_level, class_name) VALUES (5, 'A');
INSERT INTO classes (grade_level, class_name) VALUES (5, 'B');
INSERT INTO classes (grade_level, class_name) VALUES (5, 'C');
INSERT INTO classes (grade_level, class_name) VALUES (5, 'D');
INSERT INTO classes (grade_level, class_name) VALUES (5, 'F');

INSERT INTO classes (grade_level, class_name) VALUES (6, 'A');
INSERT INTO classes (grade_level, class_name) VALUES (6, 'B');
INSERT INTO classes (grade_level, class_name) VALUES (6, 'C');
INSERT INTO classes (grade_level, class_name) VALUES (6, 'D');
INSERT INTO classes (grade_level, class_name) VALUES (6, 'F');

INSERT INTO classes (grade_level, class_name) VALUES (7, 'A');
INSERT INTO classes (grade_level, class_name) VALUES (7, 'B');
INSERT INTO classes (grade_level, class_name) VALUES (7, 'C');
INSERT INTO classes (grade_level, class_name) VALUES (7, 'D');
INSERT INTO classes (grade_level, class_name) VALUES (7, 'F');

INSERT INTO classes (grade_level, class_name) VALUES (8, 'A');
INSERT INTO classes (grade_level, class_name) VALUES (8, 'B');
INSERT INTO classes (grade_level, class_name) VALUES (8, 'C');
INSERT INTO classes (grade_level, class_name) VALUES (8, 'D');
INSERT INTO classes (grade_level, class_name) VALUES (8, 'F');

INSERT INTO classes (grade_level, class_name) VALUES (9, 'A');
INSERT INTO classes (grade_level, class_name) VALUES (9, 'B');
INSERT INTO classes (grade_level, class_name) VALUES (9, 'C');
INSERT INTO classes (grade_level, class_name) VALUES (9, 'D');
INSERT INTO classes (grade_level, class_name) VALUES (9, 'F');

INSERT INTO classes (grade_level, class_name) VALUES (10, 'A');
INSERT INTO classes (grade_level, class_name) VALUES (10, 'B');
INSERT INTO classes (grade_level, class_name) VALUES (10, 'C');
INSERT INTO classes (grade_level, class_name) VALUES (10, 'D');
INSERT INTO classes (grade_level, class_name) VALUES (10, 'F');

INSERT INTO classes (grade_level, class_name) VALUES (11, 'A');
INSERT INTO classes (grade_level, class_name) VALUES (11, 'B');
INSERT INTO classes (grade_level, class_name) VALUES (11, 'C');
INSERT INTO classes (grade_level, class_name) VALUES (11, 'D');
INSERT INTO classes (grade_level, class_name) VALUES (11, 'F');

INSERT INTO classes (grade_level, class_name) VALUES (12, 'A');
INSERT INTO classes (grade_level, class_name) VALUES (12, 'B');
INSERT INTO classes (grade_level, class_name) VALUES (12, 'C');
INSERT INTO classes (grade_level, class_name) VALUES (12, 'D');
INSERT INTO classes (grade_level, class_name) VALUES (12, 'F');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM classes;
-- +goose StatementEnd
