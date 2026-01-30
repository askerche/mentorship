-- +goose Up
-- +goose StatementBegin
ALTER TABLE links DROP COLUMN redirects_count;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
