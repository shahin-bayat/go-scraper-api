-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE users
(
    id             SERIAL PRIMARY KEY,
    email          VARCHAR(255) UNIQUE NOT NULL,
    given_name     VARCHAR(255),
    family_name    VARCHAR(255),
    name           VARCHAR(255),
    locale         VARCHAR(255),
    avatar_url     VARCHAR(255),
    verified_email BOOLEAN,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE users;
-- +goose StatementEnd
