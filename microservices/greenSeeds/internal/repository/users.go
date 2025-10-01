package repository

import "github.com/jmoiron/sqlx"

type IUsersRepository interface{}

type usersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *usersRepository {
	return &usersRepository{
		db: db,
	}
}

// func (r *usersRepository) AddUser(user models.User) (bool, error) {
// 	tx, err := r.db.Beginx()
// 	if err != nil {
// 		return false, err
// 	}

// 	committed := false
// 	defer func() {
// 		if !committed {
// 			if err := tx.Rollback(); err != nil {
// 				return
// 			}
// 		}
// 	}()

// 	const insertUserQuery = `
// INSERT INTO itops.users (
// 	full_name,
// 	password_hash
// ) VALUES (
// 	:full_name,
// 	:password_hash
// ) RETURNING id;`

// 	rows, err := tx.NamedQuery(insertUserQuery, user)
// 	if err != nil {
// 		return false, err
// 	}

// 	var uuid string
// 	if rows.Next() {
// 		if err := rows.Scan(&uuid); err != nil {
// 			return false, err
// 		}
// 	} else {
// 		// пользователь уже существует
// 		return false, nil
// 	}
// 	rows.Close()

// 	const insertRolesQuery = `
// 	INSERT INTO itops.user_roles (
// 		user_id,
// 		role_name
// 	) VALUES (
// 		:id,
// 		:role_name
// 	) ON CONFLICT (user_id, role_name) DO NOTHING;`

// 	userRoles := models.User{
// 		UUID: &uuid,
// 		Role: user.Role,
// 	}

// 	result, err := tx.NamedExec(insertRolesQuery, userRoles)
// 	if err != nil {
// 		return false, err
// 	}

// 	if affected, err := result.RowsAffected(); err != nil {
// 		return false, err
// 	} else if affected == 0 {
// 		return false, nil
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return false, err
// 	}
// 	committed = true

// 	return true, nil
// }

// func (r *usersRepository) CheckUserByFullName(fullName string) (models.User, error) {
// 	const query = "SELECT id, full_name, password_hash FROM itops.users WHERE full_name = $1;"

// 	var user models.User
// 	if err := r.db.Get(&user, query, fullName); err != nil {
// 		return models.User{}, err
// 	}

// 	return user, nil
// }

// func (r *usersRepository) CheckUserByUuid(uuid string) (models.User, error) {
// 	const query = `
// SELECT
//     itops.users.id,
//     itops.users.full_name,
//     itops.users.password_hash,
//     itops.user_roles.role_name
// FROM itops.users
// LEFT JOIN itops.user_roles ON itops.users.id = itops.user_roles.user_id
// WHERE itops.users.id = $1;`

// 	var user models.User
// 	if err := r.db.Get(&user, query, uuid); err != nil {
// 		return models.User{}, err
// 	}

// 	return user, nil
// }

// func (r *usersRepository) CheckRolesById(uuid string) (string, error) {
// 	const query = `
// 	SELECT name FROM itops.roles
// 	JOIN itops.user_roles ON itops.roles.name = itops.user_roles.role_name
// 	JOIN itops.users ON itops.users.id = itops.user_roles.user_id
// 	WHERE itops.users.id = $1;`

// 	var role string
// 	if err := r.db.Get(&role, query, uuid); err != nil {
// 		return "", err
// 	}

// 	return role, nil
// }

// func (r *usersRepository) CheckAllUsers() ([]models.User, error) {
// 	const query = `
// SELECT
//     itops.users.id,
//     itops.users.full_name,
//     itops.roles.name AS role_name
// FROM itops.users
// JOIN itops.user_roles ON itops.users.id = itops.user_roles.user_id
// JOIN itops.roles ON itops.user_roles.role_name = itops.roles.name;`

// 	var users []models.User
// 	if err := r.db.Select(&users, query); err != nil {
// 		return nil, err
// 	}

// 	return users, nil
// }

// func (r *usersRepository) UpdateRole(role models.UpdateRole) (models.UpdateRole, error) {
// 	tx, err := r.db.Beginx()
// 	if err != nil {
// 		return models.UpdateRole{}, err
// 	}

// 	committed := false
// 	defer func() {
// 		if !committed {
// 			if err := tx.Rollback(); err != nil {
// 				return
// 			}
// 		}
// 	}()

// 	const query = `
// 	UPDATE itops.user_roles
// 	SET role_name = :role_name
// 	WHERE user_id = :id;`

// 	if _, err := tx.NamedExec(query, role); err != nil {
// 		return models.UpdateRole{}, err
// 	}

// 	const query2 = `
// 	SELECT role_name
// 	FROM itops.user_roles
// 	WHERE user_id = $1;`

// 	var role_name string
// 	if err := tx.Get(&role_name, query2, role.UUID); err != nil {
// 		return models.UpdateRole{}, err
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return models.UpdateRole{}, err
// 	}
// 	committed = true

// 	return models.UpdateRole{
// 		UUID: role.UUID,
// 		Role: &role_name,
// 	}, nil
// }

// func (r *usersRepository) UpdatePassword(password models.UpdatePassword) (bool, error) {
// 	const query = `
// 	UPDATE itops.users
// 	SET password_hash = :password_hash
// 	WHERE id = :id;`

// 	result, err := r.db.NamedExec(query, password)
// 	if err != nil {
// 		return false, err
// 	}

// 	affected, err := result.RowsAffected()
// 	if err != nil {
// 		return false, err
// 	}

// 	return affected == 1, nil
// }

// func (r *usersRepository) RemoveUser(uuid string) (bool, error) {
// 	tx, err := r.db.Beginx()
// 	if err != nil {
// 		return false, err
// 	}

// 	committed := false
// 	defer func() {
// 		if !committed {
// 			if err := tx.Rollback(); err != nil {
// 				return
// 			}
// 		}
// 	}()

// 	const query = `
// 	DELETE FROM itops.user_roles
// 	WHERE user_id = $1;`

// 	if _, err := tx.Exec(query, uuid); err != nil {
// 		return false, err
// 	}

// 	const query2 = `
// 	DELETE FROM itops.users
// 	WHERE id = $1;`

// 	if _, err := tx.Exec(query2, uuid); err != nil {
// 		return false, err
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return false, err
// 	}
// 	committed = true

// 	return true, nil
// }
