
CREATE TABLE IF NOT EXISTS category (
  category_id SERIAL PRIMARY KEY,
  title VARCHAR(100) NOT NULL UNIQUE,
  slug VARCHAR(100) NOT NULL UNIQUE,
  short_title VARCHAR(60) NOT NULL,
  img_path TEXT,
  parent_category_id INT REFERENCES category (category_id) ON DELETE CASCADE CHECK (category_id != parent_category_id)
);



INSERT INTO category (title, slug, short_title) VALUES ('Мужская одежда, обувь и аксессуары', 'men', 'Мужчинам'); --1

INSERT INTO category (title, slug, short_title) VALUES ('Женская одежда, обувь и аксессуары', 'women', 'Женщинам'); -- 2

INSERT INTO category (title, slug, short_title) VALUES ('Детская одежда, обувь и аксессуары', 'children', 'Детям'); -- 3

INSERT INTO category (title, slug, short_title, parent_category_id) VALUES ('Мужская Одежда', 'muzhskaya-odezhda', 'Одежда', 1); -- 4

INSERT INTO category (title, slug, short_title, parent_category_id) VALUES ('Мужская верхняя одежда', 'muzhskaya-verkhnyaya-odezhda', 'Верхняя одежда', 4); --5

INSERT INTO category (title, slug, short_title, parent_category_id) VALUES ('Мужские бомберы', 'muzhbombery', 'Бомберы', 5); --6