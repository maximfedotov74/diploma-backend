CREATE TABLE category_type (
  category_type_id SERIAL PRIMARY KEY,
  title VARCHAR(40) NOT NULL UNIQUE
);

CREATE TABLE category (
  category_id SERIAL PRIMARY KEY,
  title VARCHAR(40) NOT NULL UNIQUE,
  img_path VARCHAR(255),
  parent_category_id INT REFERENCES category (category_id) ON DELETE CASCADE CHECK (category_id != parent_category_id)
);

CREATE TABLE category_type_category (
  category_type_category_id SERIAL PRIMARY KEY,
  category_type_id INT REFERENCES category_type (category_type_id) ON DELETE CASCADE NOT NULL,
  category_id INT REFERENCES category (category_id) ON DELETE CASCADE NOT NULL
);
