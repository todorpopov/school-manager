-- +goose Up
-- +goose StatementBegin
INSERT INTO schools (school_name, school_address) VALUES ('GPAE', 'Burgas, Bulgaria');
INSERT INTO schools (school_name, school_address) VALUES ('SMG', 'Sofia, Bulgaria');
INSERT INTO schools (school_name, school_address) VALUES ('GPNE', 'Burgas, Bulgaria');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM schools;
-- +goose StatementEnd
