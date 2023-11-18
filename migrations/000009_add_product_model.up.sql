CREATE TABLE IF NOT EXISTS product_model (
  product_model_id SERIAL PRIMARY KEY,
  price float4 NOT NULL,
  discount INT2,
  product_id INT REFERENCES product (product_id) ON DELETE CASCADE NOT NULL
);



-- //todo add feedback to model 