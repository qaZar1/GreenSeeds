package models

type Users struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	FullName string `json:"full_name" db:"full_name"`
	IsAdmin  bool   `json:"is_admin" db:"is_admin"`
}
