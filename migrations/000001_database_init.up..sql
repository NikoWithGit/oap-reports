CREATE TABLE completed_orders (
    id VARCHAR(36) PRIMARY KEY NOT NULL, 
    total REAL,
    date TIMESTAMP
);

CREATE TABLE products_in_orders (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    uuid VARCHAR(36),
    num INTEGER,
    price_per_one REAL, 
    order_id VARCHAR(36)
);