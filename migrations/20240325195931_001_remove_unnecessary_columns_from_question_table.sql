-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE questions DROP COLUMN text;
ALTER TABLE questions DROP COLUMN image_path;
ALTER TABLE questions DROP COLUMN is_fetched;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE questions ADD COLUMN text VARCHAR(255);
ALTER TABLE questions ADD COLUMN image_path VARCHAR(255);
ALTER TABLE questions ADD COLUMN is_fetched BOOLEAN;
-- +goose StatementEnd
