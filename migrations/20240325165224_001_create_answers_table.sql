-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE
  IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
		question_id INTEGER NOT NULL REFERENCES questions(id),
		answer VARCHAR(255) NOT NULL,
		is_correct BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP 
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS answers;
-- +goose StatementEnd
