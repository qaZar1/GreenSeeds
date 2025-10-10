package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IUsersRepository interface {
	AddUser(user models.User) (bool, error)
	CheckUserByUsername(username string) (models.User, error)
	CheckAllUsers() ([]models.User, error)
	Update(user models.User) (bool, error)
	Delete(username string) (bool, error)
	UpdatePassword(user models.UpdatePassword) (bool, error)
	CheckUserByUsernameWithPwd(username string) (models.User, error)
}

type usersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *usersRepository {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) AddUser(user models.User) (bool, error) {
	const query = `
INSERT INTO green_seeds.users (
	username,
	password,
	full_name,
	is_admin
) VALUES (
	:username,
	:password,
	:full_name,
	:is_admin
);`

	result, err := r.db.NamedExec(query, user)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (r *usersRepository) CheckUserByUsername(username string) (models.User, error) {
	const query = "SELECT username, full_name, is_admin FROM green_seeds.users WHERE username = $1;"

	var user models.User
	if err := r.db.Get(&user, query, username); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *usersRepository) CheckUserByUsernameWithPwd(username string) (models.User, error) {
	const query = "SELECT username, password, full_name, is_admin FROM green_seeds.users WHERE username = $1;"

	var user models.User
	if err := r.db.Get(&user, query, username); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *usersRepository) CheckAllUsers() ([]models.User, error) {
	const query = `
SELECT
    username,
	full_name,
	is_admin
FROM green_seeds.users
ORDER BY username ASC;`

	var users []models.User
	if err := r.db.Select(&users, query); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *usersRepository) Update(user models.User) (bool, error) {
	const query = `
UPDATE green_seeds.users
SET full_name = COALESCE(:full_name, full_name),
	is_admin = COALESCE(:is_admin, is_admin)
WHERE username = :username;
`

	result, err := r.db.NamedExec(query, user)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (r *usersRepository) Delete(username string) (bool, error) {
	const query = `
DELETE FROM green_seeds.users
WHERE username = $1;
`

	result, err := r.db.Exec(query, username)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (r *usersRepository) UpdatePassword(user models.UpdatePassword) (bool, error) {
	const query = `
UPDATE green_seeds.users
SET password = COALESCE(:new_password, password)
WHERE username = :username;
`

	result, err := r.db.NamedExec(query, user)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}
