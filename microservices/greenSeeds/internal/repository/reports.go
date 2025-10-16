package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IReportsRepository interface {
	AddReports(reports models.Reports) (bool, error)
	GetReports() ([]models.Reports, error)
	UpdateReports(reports models.Reports) (bool, error)
	DeleteReports(shift int, number int, receipt int) (bool, error)
	GetReportsById(id int) (models.Reports, error)
}

type reportsRepository struct {
	db *sqlx.DB
}

func NewReportsRepository(db *sqlx.DB) *reportsRepository {
	return &reportsRepository{
		db: db,
	}
}

func (rep *reportsRepository) AddReports(reports models.Reports) (bool, error) {
	const query = `
INSERT INTO green_seeds.reports (
	shift,
	number,
	receipt,
	turn,
	dt,
	success,
	error,
	solution,
	mark)
VALUES (
	:shift,
	:number,
	:receipt,
	:turn,
	:dt,
	:success,
	:error,
	:solution,
	:mark)`

	result, err := rep.db.NamedExec(query, reports)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (rep *reportsRepository) GetReports() ([]models.Reports, error) {
	const query = `
SELECT
	id,
	shift,
	number,
	receipt,
	turn,
	dt,
	success,
	error,
	solution,
	mark
FROM green_seeds.reports`

	var reports []models.Reports
	if err := rep.db.Select(&reports, query); err != nil {
		return nil, err
	}

	return reports, nil
}

func (rep *reportsRepository) UpdateReports(reports models.Reports) (bool, error) {
	const query = `
UPDATE green_seeds.reports
SET turn = :turn,
	dt = :dt,
	success = :success,
	error = :error,
	solution = :solution,
	mark = :mark
WHERE
	shift = :shift AND
	number = :number AND
	receipt = :receipt
`

	result, err := rep.db.NamedExec(query, reports)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (rep *reportsRepository) DeleteReports(shift int, number int, receipt int) (bool, error) {
	const query = `
DELETE FROM green_seeds.reports
WHERE shift = $1 AND number = $2 AND receipt = $3`

	result, err := rep.db.Exec(query, shift, number, receipt)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (rep *reportsRepository) GetReportsById(id int) (models.Reports, error) {
	const query = `
SELECT
	id,
	shift,
	number,
	receipt,
	turn,
	dt,
	success,
	error,
	solution,
	mark
FROM green_seeds.reports
WHERE id = $1`

	var report models.Reports
	if err := rep.db.Get(&report, query, id); err != nil {
		return models.Reports{}, err
	}

	return report, nil
}
