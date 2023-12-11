CREATE TABLE IF NOT EXISTS product_model (
  product_model_id SERIAL PRIMARY KEY,
  price int NOT NULL,
  discount INT2,
  main_image_path TEXT NOT NULL, 
  product_id INT REFERENCES product (product_id) ON DELETE CASCADE NOT NULL
);



-- //todo add feedback to model 