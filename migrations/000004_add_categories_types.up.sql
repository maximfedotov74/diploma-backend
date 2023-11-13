
CREATE TABLE IF NOT EXISTS category (
  category_id SERIAL PRIMARY KEY,
  title VARCHAR(100) NOT NULL UNIQUE,
  slug VARCHAR(100) NOT NULL UNIQUE,
  short_title VARCHAR(60) NOT NULL,
  img_path VARCHAR(255),
  parent_category_id INT REFERENCES category (category_id) ON DELETE CASCADE CHECK (category_id != parent_category_id)
);




