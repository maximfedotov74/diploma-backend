CREATE TABLE category_type (
  category_type_id SERIAL PRIMARY KEY,
  title VARCHAR(80) NOT NULL UNIQUE
);

CREATE TABLE category (
  category_id SERIAL PRIMARY KEY,
  title VARCHAR(80) NOT NULL UNIQUE,
  img_path VARCHAR(255),
  parent_category_id INT REFERENCES category (category_id) ON DELETE CASCADE CHECK (category_id != parent_category_id)
);

CREATE TABLE category_type_category (
  category_type_category_id SERIAL PRIMARY KEY,
  category_type_id INT REFERENCES category_type (category_type_id) ON DELETE CASCADE NOT NULL,
  category_id INT REFERENCES category (category_id) ON DELETE CASCADE NOT NULL
);


select parent.category_id as parent_id, parent.title as parent_title, parent.img_path as parent_img_path,
parent.parent_category_id as parent_parent_id,
child.category_id as child_id, child.title as child_title, child.img_path as child_img_path,
child.parent_category_id as child_parent_id
from category as parent
left join category as child on parent.category_id = child.parent_category_id where parent.title = $1;