CREATE TABLE public.user
(
  user_id SERIAL PRIMARY KEY,
  email VARCHAR(129) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL
);
