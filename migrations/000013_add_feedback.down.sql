ALTER TABLE feedback DROP CONSTRAINT feedback_user_id_product_model_id_unique;
DROP TABLE IF EXISTS feedback CASCADE;
