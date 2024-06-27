-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE failed_questions
(
    id          SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL REFERENCES questions (id),
    user_id     INTEGER NOT NULL REFERENCES users (id),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP,
    UNIQUE (question_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE failed_questions;
-- +goose StatementEnd
