package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IShiftsRepository interface {
	AddShifts(shifts *models.Shifts) (bool, error)
	GetShifts() ([]models.Shifts, error)
	UpdateShifts(shifts *models.Shifts) (bool, error)
	DeleteShifts(shift int) (bool, error)
}

type shiftsRepository struct {
	db *sqlx.DB
}

func NewShiftsRepository(db *sqlx.DB) *shiftsRepository {
	return &shiftsRepository{
		db: db,
	}
}

func (sh *shiftsRepository) AddShifts(shifts *models.Shifts) (bool, error) {
	const query = `
INSERT INTO green_seeds.shifts (
	shift,
	dt,
	username
)
VALUES (
	:shift,
	:dt,
	:username
)`

	result, err := sh.db.NamedExec(query, shifts)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (sh *shiftsRepository) GetShifts() ([]models.Shifts, error) {
	const query = `
SELECT shift, dt, username
FROM green_seeds.shifts`

	var shifts []models.Shifts
	if err := sh.db.Select(&shifts, query); err != nil {
		return nil, err
	}

	return shifts, nil
}

func (sh *shiftsRepository) UpdateShifts(shifts *models.Shifts) (bool, error) {
	const query = `
UPDATE green_seeds.shifts
SET dt = :dt,
	username = :username
WHERE shift = :shift`

	result, err := sh.db.NamedExec(query, shifts)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (sh *shiftsRepository) DeleteShifts(shift int) (bool, error) {
	const query = `
DELETE FROM green_seeds.shifts
WHERE shift = $1`

	result, err := sh.db.Exec(query, shift)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}
