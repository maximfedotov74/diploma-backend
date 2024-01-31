ALTER TABLE order_delivery_point DROP CONSTRAINT delivery_point_id_order_id_unique;
DROP TABLE IF EXISTS public.order CASCADE;
DROP TABLE IF EXISTS order_model CASCADE;
DROP TABLE IF EXISTS order_activation CASCADE;
DROP TABLE IF EXISTS delivery_point CASCADE;
DROP TABLE IF EXISTS order_delivery_point CASCADE;
