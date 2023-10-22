package repository

const (
	findRoleByTitle = `SELECT role_id, title FROM public.role WHERE public.role.title = $1;`
	addRoleToUser   = "INSERT INTO public.user_role (user_id, role_id) VALUES ($1, $2);"
)
