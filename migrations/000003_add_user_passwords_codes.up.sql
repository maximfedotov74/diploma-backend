CREATE TABLE IF NOT EXISTS change_password_code (
  change_password_code_id SERIAL PRIMARY KEY,
  code VARCHAR(6) NOT NULL DEFAULT LPAD(FLOOR(RANDOM() * 1000000)::VARCHAR, 6, '0'),
  end_time timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '30 minutes',
  user_id INT UNIQUE REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL
);