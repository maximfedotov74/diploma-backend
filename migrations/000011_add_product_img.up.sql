CREATE TABLE product_model_img (
  product_img_id SERIAL PRIMARY KEY,
  img_path TEXT NOT NULL,
  product_model_id INT REFERENCES product_model (product_model_id) ON DELETE CASCADE NOT NULL
);