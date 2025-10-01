package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	UsrRepo IUsersRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UsrRepo: NewUsersRepository(db),
	}
}
