CREATE TABLE IF NOT EXISTS cart (
  cart_id SERIAL PRIMARY KEY,
  user_id INT REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL,
  model_size_id INT REFERENCES model_sizes (model_size_id) ON DELETE CASCADE NOT NULL,
  quantity int NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS wish (
  wish_id SERIAL PRIMARY KEY,
  user_id INT REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL,
  product_model_id INT REFERENCES product_model (product_model_id) ON DELETE CASCADE NOT NULL
);
