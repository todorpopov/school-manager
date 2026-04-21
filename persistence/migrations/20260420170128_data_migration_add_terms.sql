-- +goose Up
-- +goose StatementBegin
INSERT INTO terms (name) VALUES ('Fall 2026');
INSERT INTO terms (name) VALUES ('Spring 2026');
INSERT INTO terms (name) VALUES ('Fall 2027');
INSERT INTO terms (name) VALUES ('Spring 2027');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM terms;
-- +goose StatementEnd
