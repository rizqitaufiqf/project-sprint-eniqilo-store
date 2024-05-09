CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(100) NOT NULL PRIMARY KEY,
    customer_id VARCHAR(100) NOT NULL,
    product_details jsonb NOT NULL,
    paid INT NOT NULL,
    change INT NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE NO ACTION,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);