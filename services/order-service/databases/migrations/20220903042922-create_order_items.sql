
-- +migrate Up
CREATE TABLE order_items
(
    order_item_id INTEGER PRIMARY KEY,
    order_id INTEGER NULL,
    sku VARCHAR(50) NOT NULL,
    qty INTEGER NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL
);

-- +migrate Down
DROP TABLE order_items;
