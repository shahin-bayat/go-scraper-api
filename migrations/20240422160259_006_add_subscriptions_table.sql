-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    price INTEGER NOT NULL,
    currency VARCHAR(3) NOT NULL,
    interval INTERVAL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CHECK (price > 0),
    CHECK (currency ~ '^[A-Z]{3}$')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE subscriptions;
-- +goose StatementEnd
