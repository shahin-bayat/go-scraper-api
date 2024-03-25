-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE
  IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
		category VARCHAR(50) NOT NULL,
		category_key VARCHAR(50) UNIQUE NOT NULL,
		created_at TIMEStAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMEStAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP 
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
