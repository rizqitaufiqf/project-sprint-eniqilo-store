CREATE INDEX idx_staffs_phone on staffs (phone_number);
CREATE INDEX idx_products_id on products (id);
CREATE INDEX idx_customers_id on customers (id);
CREATE INDEX idx_customers_phone on customers (phone_number);

CREATE INDEX IF NOT EXISTS products_name
	ON products USING HASH(lower(name));
-- CREATE INDEX IF NOT EXISTS products_user_id
-- 	ON products (user_id);
CREATE INDEX IF NOT EXISTS products_sku
	ON products USING HASH(sku);
CREATE INDEX IF NOT EXISTS products_category
	ON products (category);
CREATE INDEX IF NOT EXISTS products_is_available
	ON products (is_available);
CREATE INDEX IF NOT EXISTS products_in_stock
	ON products (stock) WHERE stock > 0;
CREATE INDEX IF NOT EXISTS products_not_in_stock
	ON products (stock) WHERE stock = 0;
CREATE INDEX IF NOT EXISTS products_created_at_desc
	ON products(created_at DESC);
CREATE INDEX IF NOT EXISTS products_created_at_asc
	ON products(created_at ASC);
CREATE INDEX IF NOT EXISTS products_price_desc
	ON products(price DESC);
CREATE INDEX IF NOT EXISTS products_price_asc
	ON products(price ASC);