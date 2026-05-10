-- +goose Up
-- +goose StatementBegin
ALTER TABLE order_items ADD COLUMN quantity INT NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
