
-- +migrate Up
CREATE TABLE customers
(
    customer_id INTEGER PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password VARCHAR(255) NULL,
    status SMALLINT DEFAULT '10'
);

COMMENT ON COLUMN customers.status IS '0/10';

-- +migrate Down
DROP TABLE customers;
