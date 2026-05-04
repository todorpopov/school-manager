-- +goose Up
-- +goose StatementBegin
INSERT INTO terms (name, start_date, end_date) VALUES ('Fall 2026', '2026-09-01', '2027-01-15');
INSERT INTO terms (name, start_date, end_date) VALUES ('Spring 2026', '2026-01-15', '2026-06-01');
INSERT INTO terms (name, start_date, end_date) VALUES ('Fall 2027', '2027-09-01', '2028-01-15');
INSERT INTO terms (name, start_date, end_date) VALUES ('Spring 2027', '2027-01-15', '2027-06-01');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM terms;
-- +goose StatementEnd
