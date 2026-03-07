-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (role_name) VALUES ('ADMIN');
INSERT INTO roles (role_name) VALUES ('USER');
INSERT INTO roles (role_name) VALUES ('EMPLOYEE');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM roles WHERE role_name IN ('ADMIN', 'USER', 'EMPLOYEE');
-- +goose StatementEnd
