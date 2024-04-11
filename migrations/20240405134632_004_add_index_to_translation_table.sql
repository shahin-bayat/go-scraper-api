-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE INDEX idx_translations_refer_id_type_lang ON translations (refer_id, type, lang);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX IF EXISTS idx_translations_refer_id_type_lang;
-- +goose StatementEnd
