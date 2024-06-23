-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE users DROP CONSTRAINT users_stripe_customer_id_key;
ALTER TABLE users ADD UNIQUE (id, stripe_customer_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE users ADD UNIQUE (stripe_customer_id);
ALTER TABLE users DROP CONSTRAINT users_id_stripe_customer_id_key;
-- +goose StatementEnd
