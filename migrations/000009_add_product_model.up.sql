CREATE TABLE IF NOT EXISTS product_model (
  product_model_id SERIAL PRIMARY KEY,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  price int NOT NULL,
  slug VARCHAR(255) NOT NULL UNIQUE,
  article VARCHAR(12) NOT NULL UNIQUE,
  discount INT2,
  main_image_path TEXT NOT NULL, 
  product_id INT REFERENCES product (product_id) ON DELETE CASCADE NOT NULL
);



-- //todo add feedback to model 
