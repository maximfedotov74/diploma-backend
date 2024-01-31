CREATE TABLE IF NOT EXISTS brand (
  brand_id SERIAL PRIMARY KEY,
  title VARCHAR(100) NOT NULL UNIQUE,
  description TEXT,
  img_path TEXT
);

