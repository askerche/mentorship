-- +goose Up
-- +goose StatementBegin
ALTER TABLE products ADD COLUMN image_file_id VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products DROP COLUMN image_file_id;
-- +goose StatementEnd
