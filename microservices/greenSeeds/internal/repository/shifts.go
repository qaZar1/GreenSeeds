package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IShiftsRepository interface {
	AddShifts(shifts models.Shifts) (models.Shifts, error)
	GetShifts() ([]models.Shifts, error)
	UpdateShifts(shifts models.Shifts) (models.Shifts, error)
	DeleteShifts(shift int) (bool, error)
	GetShiftsByShift(shift int) (models.Shifts, error)
	GetShiftsWithoutUser() ([]models.Shifts, error)
	GetShiftsByUsername(username string) ([]models.Shifts, error)
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
WHERE dt >= (CURRENT_DATE AT TIME ZONE 'UTC+5') - INTERVAL '7 days' AND deleted_at IS NULL
ORDER BY shift ASC;`

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
WHERE shift = :shift AND deleted_at IS NULL
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
UPDATE green_seeds.shifts
SET deleted_at = $1
WHERE shift = $2 AND deleted_at IS NULL;`

	result, err := sh.db.Exec(query, time.Now(), shift)
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
WHERE shift = $1 AND deleted_at IS NULL;`

	var shifts models.Shifts
	if err := sh.db.Get(&shifts, query, shift); err != nil {
		return models.Shifts{}, err
	}

	return shifts, nil
}

func (sh *shiftsRepository) GetShiftsWithoutUser() ([]models.Shifts, error) {
	const query = `
SELECT shift, dt, username
FROM green_seeds.shifts
WHERE DATE(dt) = CURRENT_DATE AND username IS NULL AND deleted_at IS NULL
ORDER BY shift ASC;`

	var shifts []models.Shifts
	if err := sh.db.Select(&shifts, query); err != nil {
		return nil, err
	}

	return shifts, nil
}

func (sh *shiftsRepository) GetShiftsByUsername(username string) ([]models.Shifts, error) {
	const query = `
SELECT shift, dt, username
FROM green_seeds.shifts
WHERE DATE(dt) = CURRENT_DATE AND username = $1 AND deleted_at IS NULL;`

	var shifts []models.Shifts
	if err := sh.db.Select(&shifts, query, username); err != nil {
		return nil, err
	}

	return shifts, nil
}
