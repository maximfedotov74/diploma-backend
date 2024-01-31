CREATE TABLE IF NOT EXISTS option (
  option_id SERIAL PRIMARY KEY,
  title VARCHAR(100) NOT NULL,
  slug VARCHAR(150) NOT NULL UNIQUE,
  for_catalog boolean NOT NULL DEFAULT true
);


CREATE TABLE IF NOT EXISTS option_value (
  option_value_id SERIAL PRIMARY KEY,
  value VARCHAR(150) NOT NULL,
  info VARCHAR(150),
  option_id INT REFERENCES option (option_id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS product_model_option (
  product_model_option_id SERIAL PRIMARY KEY,
  product_model_id INT REFERENCES product_model (product_model_id) ON DELETE CASCADE NOT NULL,
  option_id INT REFERENCES option (option_id) ON DELETE CASCADE NOT NULL,
  option_value_id INT REFERENCES option_value (option_value_id) ON DELETE CASCADE NOT NULL
);


