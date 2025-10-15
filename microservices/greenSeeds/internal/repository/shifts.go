package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IShiftsRepository interface {
	AddShifts(shifts models.Shifts) (models.Shifts, error)
	GetShifts() ([]models.Shifts, error)
	UpdateShifts(shifts models.Shifts) (models.Shifts, error)
	DeleteShifts(shift int) (bool, error)
	GetShiftsByShift(shift int) (models.Shifts, error)
}

type shiftsRepository struct {
	db *sqlx.DB
}

func NewShiftsRepository(db *sqlx.DB) *shiftsRepository {
	return &shiftsRepository{
		db: db,
	}
}

func (sh *shiftsRepository) AddShifts(shifts models.Shifts) (models.Shifts, error) {
	const query = `
INSERT INTO green_seeds.shifts (
	dt
)
VALUES (
	:dt
)
RETURNING shift, dt`

	rows, err := sh.db.NamedQuery(query, shifts)
	if err != nil {
		return models.Shifts{}, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&shifts.Shift, &shifts.Dt); err != nil {
			return models.Shifts{}, err
		}
	}

	return shifts, nil
}

func (sh *shiftsRepository) GetShifts() ([]models.Shifts, error) {
	const query = `
SELECT shift, dt, username
FROM green_seeds.shifts
WHERE dt >= CURRENT_DATE
ORDER BY shift ASC`

	var shifts []models.Shifts
	if err := sh.db.Select(&shifts, query); err != nil {
		return nil, err
	}

	return shifts, nil
}

func (sh *shiftsRepository) UpdateShifts(shifts models.Shifts) (models.Shifts, error) {
	const query = `
UPDATE green_seeds.shifts
SET dt = COALESCE(:dt, dt),
	username = COALESCE(:username, username)
WHERE shift = :shift
RETURNING shift, dt`

	rows, err := sh.db.NamedQuery(query, shifts)
	if err != nil {
		return models.Shifts{}, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&shifts.Shift, &shifts.Dt); err != nil {
			return models.Shifts{}, err
		}
	}

	return shifts, nil
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

func (sh *shiftsRepository) GetShiftsByShift(shift int) (models.Shifts, error) {
	const query = `
SELECT shift, dt, username
FROM green_seeds.shifts
WHERE shift = $1`

	var shifts models.Shifts
	if err := sh.db.Get(&shifts, query, shift); err != nil {
		return models.Shifts{}, err
	}

	return shifts, nil
}
