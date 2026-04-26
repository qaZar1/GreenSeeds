package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ICalibrationsRepository interface {
	Save(tx *sqlx.Tx, step string) (bool, error)
	GetMax() (float64, error)
	Delete(tx *sqlx.Tx) (bool, error)
	TxUpsert(step float64) error
}

type calibrationsRepository struct {
	db *sqlx.DB
}

func NewCalibrationsRepository(db *sqlx.DB) *calibrationsRepository {
	return &calibrationsRepository{
		db: db,
	}
}

func (cal *calibrationsRepository) Save(tx *sqlx.Tx, step string) (bool, error) {
	const query = `
INSERT INTO green_seeds.device_settings (
	key,
	value
)
VALUES (
	'step',
	$1
);`

	rows, err := tx.Exec(query, step)
	if err != nil {
		return false, err
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected == 1, nil
}

func (cal *calibrationsRepository) GetMax() (float64, error) {
	const query = `
SELECT value
FROM green_seeds.device_settings
WHERE key = 'step';`

	var max float64
	if err := cal.db.Get(&max, query); err != nil {
		return 0, err
	}

	return max, nil
}

func (cal *calibrationsRepository) Delete(tx *sqlx.Tx) (bool, error) {
	const query = `
DELETE FROM green_seeds.device_settings
WHERE key = 'step';`

	result, err := tx.Exec(query)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (cal *calibrationsRepository) TxUpsert(step float64) error {
	tx, err := cal.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ok, err := cal.Delete(tx)
	if err != nil {
		return err
	}

	stepStr := fmt.Sprintf("%v", step)

	ok, err = cal.Save(tx, stepStr)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("invalid to upsert data")
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
