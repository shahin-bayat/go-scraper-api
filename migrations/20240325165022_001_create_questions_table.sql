-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE
  IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    question TEXT NOT NULL,
    question_key VARCHAR(50) UNIQUE NOT NULL,
    question_number VARCHAR(50) NOT NULL,
    has_image BOOLEAN NOT NULL DEFAULT FALSE,
    file_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP 
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS questions;
-- +goose StatementEnd