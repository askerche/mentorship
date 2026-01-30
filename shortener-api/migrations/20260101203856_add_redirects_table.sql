-- +goose Up
-- +goose StatementBegin
CREATE TABLE redirects (
    id bigserial primary key,
    user_agent TEXT,
    short_link TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
