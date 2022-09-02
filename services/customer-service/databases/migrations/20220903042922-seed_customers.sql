
-- +migrate Up
INSERT INTO customers (
    customer_id,
    name,
    email,
    password,
    status
)
VALUES
    (
        1,
        'Bob Martin',
        'bobmartin@gmail.com',
        '$2a$12$Nh0BfLcDJq2Ul216GZ/kquxXJiAYTGksrZ6ubXXMLEuhpAb5MZL32',
        10
    ),
    (
        2,
        'Linus Torvalds',
        'torvalds@linux.com',
        '$2a$12$mqmKcSsMMC.DGGNWyAuQWO0f4zveYM4E6cBCRBOhTgtSQV4o3iUnS',
        10
    );

-- +migrate Down
DELETE FROM customers;
