CREATE TABLE IF NOT EXISTS orders
(
    id          SERIAL PRIMARY KEY,
    coupon_code VARCHAR(10),
    final_price NUMERIC(10, 2) NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items
(
    id         SERIAL PRIMARY KEY,
    order_id   INT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products (id),
    quantity   INT NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);
