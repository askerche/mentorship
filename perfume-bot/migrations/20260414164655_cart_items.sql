-- +goose Up
-- +goose StatementBegin
CREATE TABLE cart_items (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Жестко связываем корзину с юзером. Если юзер удалится — корзина очистится
    CONSTRAINT fk_user FOREIGN KEY (telegram_id) REFERENCES users(telegram_id) ON DELETE CASCADE,
    -- Жестко связываем корзину с товаром. Если товар удалят из магазина — он исчезнет из корзин
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    -- Защита от дублей: один юзер = одна запись на конкретный товар (дальше только увеличиваем quantity)
    UNIQUE(telegram_id, product_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_items;
-- +goose StatementEnd