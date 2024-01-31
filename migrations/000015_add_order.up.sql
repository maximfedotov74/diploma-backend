DROP TYPE IF EXISTS order_status_enum;
CREATE TYPE order_status_enum AS enum ('completed', 'waiting_for_activation', 'canceled', 'on_the_way', 'waiting_for_payment', 'paid', 'in_processing');

DROP TYPE IF EXISTS order_payment_method_enum;
CREATE TYPE order_payment_method_enum AS enum ('upon_receipt', 'online');

DROP TYPE IF EXISTS order_conditions_enum;
CREATE TYPE order_conditions_enum AS enum ('with_fitting', 'without_fitting');

CREATE TABLE IF NOT EXISTS public.order (
  order_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  delivery_date timestamp(3),
  is_activated bool DEFAULT false,
  order_status order_status_enum NOT NULL DEFAULT 'waiting_for_activation',
  order_payment_method order_payment_method_enum NOT NULL,
  conditions order_conditions_enum NOT NULL,
  products_price int NOT NULL,
  total_price int NOT NULL,
  total_discount int NOT NULL DEFAULT 0,
  promo_discount int NOT NULL DEFAULT 0,
  delivery_price int NOT NULL,
  recipient_firstname VARCHAR(255) NOT NULL,
  recipient_lastname VARCHAR(255) NOT NULL,
  recipient_phone VARCHAR(18) NOT NULL,
  user_id INT REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL
);


CREATE TABLE IF NOT EXISTS order_model (
  order_model_id SERIAL PRIMARY KEY,
  order_id UUID REFERENCES public.order (order_id) ON DELETE CASCADE NOT NULL,
  model_size_id INT REFERENCES model_sizes (model_size_id) ON DELETE CASCADE NOT NULL,
  quantity int NOT NULL,
  price int NOT NULL,
  discount int2
);


CREATE TABLE IF NOT EXISTS order_activation (
  order_activation_id SERIAL PRIMARY KEY,
  order_id UUID UNIQUE REFERENCES public.order (order_id) ON DELETE CASCADE NOT NULL,
  link UUID UNIQUE DEFAULT uuid_generate_v4(),
  end_time timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '4 hours'
);

CREATE TABLE IF NOT EXISTS delivery_point (
  delivery_point_id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  city VARCHAR(255) NOT NULL,
  address VARCHAR(255) NOT NULL,
  coords VARCHAR(255) NOT NULL,
  with_fitting bool NOT NULL,
  work_schedule VARCHAR(255) NOT NULL,
  info TEXT
);

CREATE TABLE IF NOT EXISTS order_delivery_point (
  order_delivery_point_id SERIAL PRIMARY KEY,
  delivery_point_id INT REFERENCES delivery_point (delivery_point_id) ON DELETE CASCADE NOT NULL,
  order_id UUID REFERENCES public.order (order_id) ON DELETE CASCADE NOT NULL
);

ALTER TABLE order_delivery_point ADD CONSTRAINT "delivery_point_id_order_id_unique" UNIQUE ("delivery_point_id", "order_id");
