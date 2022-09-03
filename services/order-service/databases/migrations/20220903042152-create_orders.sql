
-- +migrate Up
CREATE TABLE orders
(
    order_id INTEGER PRIMARY KEY,
    customer_id INTEGER NULL,
    order_no VARCHAR(50) NOT NULL,
    grand_total DOUBLE PRECISION NOT NULL,
    status VARCHAR(10) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL
);

-- +migrate Down
DROP TABLE orders;
