CREATE TABLE IF NOT EXISTS feedback (
  feedback_id SERIAL PRIMARY KEY,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_hidden boolean NOT NULL DEFAULT true, 
  feedback_text TEXT NOT NULL,
  rate INT2 NOT NULL,
  product_model_id INT REFERENCES product_model (product_model_id) ON DELETE CASCADE NOT NULL,
  user_id INT REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL
);

ALTER TABLE feedback ADD CONSTRAINT "feedback_user_id_product_model_id_unique" UNIQUE ("user_id", "product_model_id");
