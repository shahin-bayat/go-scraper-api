-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE questions ADD COLUMN deleted_at TIMESTAMP;
ALTER TABLE answers ADD COLUMN deleted_at TIMESTAMP;
ALTER TABLE category_questions ADD COLUMN deleted_at TIMESTAMP;
ALTER TABLE categories ADD COLUMN deleted_at TIMESTAMP;
ALTER TABLE images ADD COLUMN deleted_at TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE questions DROP COLUMN deleted_at;
ALTER TABLE answers DROP COLUMN deleted_at;
ALTER TABLE category_questions DROP COLUMN deleted_at;
ALTER TABLE categories DROP COLUMN deleted_at;
ALTER TABLE images DROP COLUMN deleted_at;
-- +goose StatementEnd
