CREATE TABLE IF NOT EXISTS session
(
  session_id SERIAL PRIMARY KEY,
  token TEXT NOT NULL UNIQUE,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_agent TEXT NOT NULL,
  user_id INT REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL
);