package models

type User struct {
	Id       *int64  `json:"id" db:"id"`
	Username string  `json:"username" db:"username"`
	Password *string `json:"password" db:"password"`
	FullName *string `json:"full_name" db:"full_name"`
	IsAdmin  *bool   `json:"is_admin" db:"is_admin"`
} // @name user

type UpdatePassword struct {
	Id          int64   `json:"id" db:"id"`
	OldPassword *string `json:"old_password"`
	NewPassword *string `json:"new_password" db:"new_password"`
} // @name UpdatePassword
