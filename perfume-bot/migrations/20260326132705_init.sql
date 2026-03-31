-- +goose Up
-- +goose StatementBegin
CREATE TABLE brands (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE
);

CREATE TABLE products(
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    brand_id BIGINT REFERENCES brands(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    price INTEGER NOT NULL
);

CREATE TABLE products_categories (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id),
    category_id BIGINT NOT NULL REFERENCES categories(id)
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
