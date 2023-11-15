CREATE TABLE product_img (
  product_img_id SERIAL PRIMARY KEY,
  img_path VARCHAR(255) NOT NULL,
  main boolean NOT NULL DEFAULT false,
  product_id INT REFERENCES product (product_id) ON DELETE CASCADE NOT NULL
);