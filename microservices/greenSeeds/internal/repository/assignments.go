package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IAssignmentsRepository interface {
	AddAssignments(assignments models.Assignments) (models.Assignments, error)
	GetAssignments() ([]models.Assignments, error)
	UpdateAssignments(assignments models.Assignments) (models.Assignments, error)
	DeleteAssignments(number int) (bool, error)
	GetAssignmentsByNumber(number int) (models.Assignments, error)
}

type assignmentsRepository struct {
	db *sqlx.DB
}

func NewAssignmentsRepository(db *sqlx.DB) *assignmentsRepository {
	return &assignmentsRepository{
		db: db,
	}
}

func (assign *assignmentsRepository) AddAssignments(assignments models.Assignments) (models.Assignments, error) {
	const query = `
INSERT INTO green_seeds.assignments (
	shift,
	number,
	receipt,
	amount
)
VALUES (
	:shift,
	:number,
	:receipt,
	:amount
)
ON CONFLICT (shift, number, receipt) DO NOTHING
RETURNING id, shift, number, receipt, amount`

	rows, err := assign.db.NamedQuery(query, assignments)
	if err != nil {
		return models.Assignments{}, err
	}
	defer rows.Close()

	var inserted models.Assignments
	if rows.Next() {
		if err = rows.StructScan(&inserted); err != nil {
			return models.Assignments{}, err
		}
	}

	return inserted, nil
}

func (assign *assignmentsRepository) GetAssignments() ([]models.Assignments, error) {
	const query = `
SELECT id, green_seeds.assignments.shift, number, receipt, amount
FROM green_seeds.assignments
JOIN green_seeds.shifts ON green_seeds.assignments.shift = green_seeds.shifts.shift
WHERE green_seeds.shifts.dt >= CURRENT_DATE
ORDER BY green_seeds.assignments.shift, green_seeds.assignments.number;`

	var assignments []models.Assignments
	if err := assign.db.Select(&assignments, query); err != nil {
		return nil, err
	}

	return assignments, nil
}

func (assign *assignmentsRepository) UpdateAssignments(assignments models.Assignments) (models.Assignments, error) {
	const query = `
UPDATE green_seeds.assignments
SET shift = :shift,
	number = :number,
	receipt = :receipt,
	amount = :amount
WHERE id = :id
RETURNING id, shift, number, receipt, amount;
`

	rows, err := assign.db.NamedQuery(query, assignments)
	if err != nil {
		return models.Assignments{}, err
	}
	defer rows.Close()

	var updated models.Assignments
	if rows.Next() {
		if err = rows.StructScan(&updated); err != nil {
			return models.Assignments{}, err
		}
	}

	return updated, nil
}

func (assign *assignmentsRepository) DeleteAssignments(id int) (bool, error) {
	const query = `
DELETE FROM green_seeds.assignments
WHERE id = $1`

	result, err := assign.db.Exec(query, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (assign *assignmentsRepository) GetAssignmentsByNumber(id int) (models.Assignments, error) {
	const query = `
SELECT id, shift, number, receipt, amount
FROM green_seeds.assignments
WHERE id = $1;`

	var assignments models.Assignments
	if err := assign.db.Get(&assignments, query, id); err != nil {
		return models.Assignments{}, err
	}

	return assignments, nil
}
