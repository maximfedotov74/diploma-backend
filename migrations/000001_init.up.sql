CREATE TABLE public.user
(
  user_id SERIAL PRIMARY KEY,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  email VARCHAR(129) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE role
(
  role_id SERIAL PRIMARY KEY,
  title VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE public.user_role
(
  user_role_id SERIAL PRIMARY KEY,
  user_id INT REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL,
  role_id INT REFERENCES public.role (role_id) ON DELETE CASCADE NOT NULL
);

CREATE TYPE auth_type AS ENUM ('credentials', 'yandex');
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.user_settings (
  user_settings_id SERIAL PRIMARY KEY,
  activation_account_link UUID DEFAULT NULL,
  is_activated boolean NOT NULL DEFAULT false,
  auth_provider auth_type NOT NULL,
  user_id INT UNIQUE REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL
);
