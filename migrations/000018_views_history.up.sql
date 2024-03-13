CREATE TABLE IF NOT EXISTS views_history (
  views_history_id SERIAL PRIMARY KEY,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id INT REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL,
  product_model_id INT REFERENCES product_model (product_model_id) ON DELETE CASCADE NOT NULL
)