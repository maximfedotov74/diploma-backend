CREATE TABLE IF NOT EXISTS public.action (
  action_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  end_date timestamp(3),
  title VARCHAR(255) NOT NULL,
  is_activated boolean NOT NULL DEFAULT false,
  img_path TEXT,
  description TEXT
);


CREATE TABLE IF NOT EXISTS action_model (
  action_model_id SERIAL PRIMARY KEY,
  action_id UUID REFERENCES public.action (action_id) ON DELETE CASCADE NOT NULL,
  product_model_id INT REFERENCES product_model (product_model_id)
);