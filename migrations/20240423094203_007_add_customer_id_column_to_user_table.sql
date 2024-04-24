-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE users ADD COLUMN stripe_customer_id varchar(50) UNIQUE DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE users DROP COLUMN stripe_customer_id;
-- +goose StatementEnd
