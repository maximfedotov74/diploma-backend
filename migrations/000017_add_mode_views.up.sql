CREATE TABLE IF NOT EXISTS product_model_views (
  product_model_views_id SERIAL PRIMARY KEY,
  ip INET NOT NULL,
  product_model_id INT REFERENCES product_model (product_model_id) ON DELETE CASCADE NOT NULL
)