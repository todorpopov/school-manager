-- +goose Up
-- +goose StatementBegin
INSERT INTO subjects (subject_name) VALUES ('Bulgarian');
INSERT INTO subjects (subject_name) VALUES ('Maths');
INSERT INTO subjects (subject_name) VALUES ('Literature');
INSERT INTO subjects (subject_name) VALUES ('Physics');
INSERT INTO subjects (subject_name) VALUES ('Chemistry');
INSERT INTO subjects (subject_name) VALUES ('Biology');
INSERT INTO subjects (subject_name) VALUES ('English');
INSERT INTO subjects (subject_name) VALUES ('History');
INSERT INTO subjects (subject_name) VALUES ('Geography');
INSERT INTO subjects (subject_name) VALUES ('Physical Education');
INSERT INTO subjects (subject_name) VALUES ('Music');
INSERT INTO subjects (subject_name) VALUES ('Information Technology');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM subjects;
-- +goose StatementEnd
