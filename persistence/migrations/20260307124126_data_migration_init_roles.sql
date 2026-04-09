-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (role_name) VALUES ('ADMIN');
INSERT INTO roles (role_name) VALUES ('DIRECTOR');
INSERT INTO roles (role_name) VALUES ('TEACHER');
INSERT INTO roles (role_name) VALUES ('STUDENT');
INSERT INTO roles (role_name) VALUES ('PARENT');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM roles WHERE role_name IN ('ADMIN', 'DIRECTOR', 'TEACHER', 'STUDENT', 'PARENT');
-- +goose StatementEnd
