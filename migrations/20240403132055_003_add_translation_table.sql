-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TYPE type AS ENUM ('question', 'answer');
CREATE TYPE lang AS ENUM ('en');
CREATE TABLE IF NOT EXISTS translations (
  id SERIAL PRIMARY KEY,
  refer_id INTEGER NOT NULL,
  type type NOT NULL,
  lang lang NOT NULL,
  translation TEXT NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TYPE type;
DROP TYPE lang;
DROP TABLE translations;
-- +goose StatementEnd
