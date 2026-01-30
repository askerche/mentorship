-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS links(
		id BIGSERIAL PRIMARY KEY,
		short_url TEXT NOT NULL,
		long_url TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
