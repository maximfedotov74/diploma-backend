CREATE TABLE IF NOT EXISTS sizes (
  size_id SERIAL PRIMARY KEY,
  size_value VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS model_sizes (
  model_size_id SERIAL PRIMARY KEY,
  product_model_id INT REFERENCES product_model (product_model_id) ON DELETE CASCADE NOT NULL,
  size_id INT REFERENCES sizes (size_id) ON DELETE CASCADE NOT NULL,
  literal_size VARCHAR(10) NOT NULL,
  in_stock INT NOT NULL DEFAULT 0
);


-- insert

-- INSERT INTO brand (title, slug) values ('Adidas', 'adidas'); -- 1
-- INSERT INTO brand (title, slug) values ('Ostin', 'ostin'); -- 2
-- INSERT INTO brand (title, slug) values ('PUMA', 'puma'); -- 3


-- INSERT INTO product (title, slug, category_id, brand_id) values ('Теплный бомбер', 'tepliy-bobmer', 6, 1);



-- INSERT INTO product_model (price, product_id, main_image_path) values (23456.99, 1, '/static/tepliy-bomber-1.webp'); -- 1
-- INSERT INTO product_model (price, product_id, main_image_path) values (56113.99, 1, '/static/tepliy-bomber-2.webp'); -- 2


-- INSERT INTO option (title,slug) values ('Материал', 'material'); -- 1
-- INSERT INTO option (title,slug) values ('Цвет', 'color'); -- 2


-- INSERT INTO option_value (value,option_id) values ('Полиэстер', 1); -- 1
-- INSERT INTO option_value (value,option_id) values ('Хлопок', 1); -- 2
-- INSERT INTO option_value (value,option_id) values ('Шерсть', 1); -- 3


-- INSERT INTO option_value (value,option_id) values ('Желтый', 2); -- 4
-- INSERT INTO option_value (value,option_id) values ('Зеленый', 2); -- 5
-- INSERT INTO option_value (value,option_id) values ('Красный', 2); -- 6
-- INSERT INTO option_value (value,option_id) values ('Натуральная кожа', 1); -- 7



-- INSERT INTO product_model_option (product_model_id,option_id,option_value_id) values (1, 1, 1);
-- INSERT INTO product_model_option (product_model_id,option_id,option_value_id) values (1, 2, 4);

-- INSERT INTO product_model_option (product_model_id,option_id,option_value_id) values (2, 1, 2);
-- INSERT INTO product_model_option (product_model_id,option_id,option_value_id) values (2, 2, 5);





-- INSERT INTO sizes (size_value) values ('44'); -- 1
-- INSERT INTO sizes (size_value) values ('46'); -- 2
-- INSERT INTO sizes (size_value) values ('48'); -- 3
-- INSERT INTO sizes (size_value) values ('50'); -- 4
-- INSERT INTO sizes (size_value) values ('40'); -- 5
-- INSERT INTO sizes (size_value) values ('41'); -- 6


-- INSERT INTO model_sizes (product_model_id, size_id, in_stock, literal_size) values (1, 1, 100, 'S');
-- INSERT INTO model_sizes (product_model_id, size_id, in_stock, literal_size) values (1, 2, 150,  'M');

-- INSERT INTO model_sizes (product_model_id, size_id, in_stock, literal_size) values (2, 3, 200, 'L');
-- INSERT INTO model_sizes (product_model_id, size_id, in_stock, literal_size) values (2, 4, 300, 'XL');


-- INSERT INTO category (title, slug, short_title, parent_category_id) VALUES ('Мужская Обувь', 'muzhskaya-obuv', 'Обувь', 1); -- 7

-- INSERT INTO category (title, slug, short_title, parent_category_id) VALUES ('Мужские ботинки', 'muzhskie-botinki', 'Ботинки', 7); --8

-- INSERT INTO category (title, slug, short_title, parent_category_id) VALUES ('Мужские высокие ботинки', 'vysokiebotinkimuj', 'Высокие ботинки', 8); --9


-- INSERT INTO product (title, slug, category_id, brand_id) values ('Ботинки', 'shoes-ostin-botinki', 9, 2);

-- INSERT INTO product_model (price, product_id, main_image_path) values (9999.99, 2, '/static/botinki-chernie.webp'); -- 1


-- INSERT INTO option_value (value,option_id) values ('Черный', 2); -- 8


-- INSERT INTO product_model_option (product_model_id,option_id,option_value_id) values (3, 1, 7);
-- INSERT INTO product_model_option (product_model_id,option_id,option_value_id) values (3, 2, 8);

-- INSERT INTO model_sizes (product_model_id, size_id, in_stock, literal_size) values (3, 5, 40, '40 RUS');
-- INSERT INTO model_sizes (product_model_id, size_id, in_stock, literal_size) values (3, 6, 60, '41 RUS');

-- insert into product_model_img (img_path,product_model_id) values ('/test1.png', 1);
-- insert into product_model_img (img_path,product_model_id) values ('/test2.png', 2);
-- insert into product_model_img (img_path,product_model_id) values ('/test3.png', 3);