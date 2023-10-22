CREATE TABLE public.user_role
(
  user_role_id SERIAL PRIMARY KEY,
  user_id INT REFERENCES public.user (user_id) ON DELETE CASCADE NOT NULL,
  role_id INT REFERENCES public.role (role_id) ON DELETE CASCADE NOT NULL
);