package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IReportsRepository interface {
	AddReports(reports models.Reports) (models.Reports, error)
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

func (rep *reportsRepository) AddReports(reports models.Reports) (models.Reports, error) {
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
	:mark)
ON CONFLICT (shift, number, receipt, turn) DO NOTHING
RETURNING id`

	rows, err := rep.db.NamedQuery(query, reports)
	if err != nil {
		return models.Reports{}, err
	}

	defer rows.Close()

	var inserted models.Reports
	if rows.Next() {
		err = rows.StructScan(&inserted)
		if err != nil {
			return models.Reports{}, err
		}
	}

	return inserted, nil
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
	COALESCE(success, FALSE) as success,
	error,
	solution,
	mark
FROM green_seeds.reports
ORDER BY shift DESC`

	var reports []models.Reports
	if err := rep.db.Select(&reports, query); err != nil {
		return nil, err
	}

	return reports, nil
}

func (rep *reportsRepository) UpdateReports(reports models.Reports) (bool, error) {
	const query = `
UPDATE green_seeds.reports
SET	dt = COALESCE(:dt, dt),
	success = COALESCE(:success, success),
	error = COALESCE(:error, error),
	solution = COALESCE(:solution, solution),
	mark = COALESCE(:mark, mark)
WHERE
	shift = :shift AND
	number = :number AND
	receipt = :receipt AND
	turn = :turn
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
	r.id,
	r.shift,
	r.number,
	r.receipt,
	r.turn,
	r.dt,
	COALESCE(r.success, FALSE) as success,
	r.error,
	r.solution,
	r.mark,
	u.full_name
FROM green_seeds.reports r
left join green_seeds.shifts s
on s.shift = r.shift
left join green_seeds.users u
on s.username = u.username
WHERE id = $1;`

	var report models.Reports
	if err := rep.db.Get(&report, query, id); err != nil {
		return models.Reports{}, err
	}

	return report, nil
}
