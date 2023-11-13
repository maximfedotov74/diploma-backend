package role

type Role struct {
	Id    int    `json:"role_id" db:"role_id" example:"1"`
	Title string `json:"title" db:"title" example:"ADMIN"`
}
