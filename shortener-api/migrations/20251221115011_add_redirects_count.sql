-- +goose Up
-- +goose StatementBegin
ALTER TABLE links ADD COLUMN redirects_count INTEGER DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
