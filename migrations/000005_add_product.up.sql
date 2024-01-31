CREATE TABLE IF NOT EXISTS product (
  product_id SERIAL PRIMARY KEY,
  title VARCHAR(200) NOT NULL,
  description TEXT,
  category_id INT REFERENCES public.category (category_id) ON DELETE CASCADE NOT NULL
);