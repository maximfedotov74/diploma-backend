CREATE TYPE auth_type AS ENUM ('credentials', 'yandex');
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_settings (
  user_settings_id SERIAL PRIMARY KEY,
  activation_account_link UUID DEFAULT NULL,
  auth_provider auth_type NOT NULL,
  user_id INT UNIQUE REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL
);