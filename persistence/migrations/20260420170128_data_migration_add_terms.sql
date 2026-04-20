-- +goose Up
-- +goose StatementBegin
INSERT INTO terms (name) VALUES ('Fall 2026');
INSERT INTO terms (name) VALUES ('Spring 2026');
INSERT INTO terms (name) VALUES ('Fall 2027');
INSERT INTO terms (name) VALUES ('Spring 2027');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM terms WHERE name IN ('Fall 2026', 'Spring 2026', 'Fall 2027', 'Spring 2027');
-- +goose StatementEnd
