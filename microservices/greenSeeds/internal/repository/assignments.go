package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IAssignmentsRepository interface {
	AddAssignments(assignments *models.Assignments) (bool, error)
	GetAssignments() ([]models.Assignments, error)
	UpdateAssignments(assignments *models.Assignments) (bool, error)
	DeleteAssignments(id int64) (bool, error)
}

type assignmentsRepository struct {
	db *sqlx.DB
}

func NewAssignmentsRepository(db *sqlx.DB) *assignmentsRepository {
	return &assignmentsRepository{
		db: db,
	}
}

func (assign *assignmentsRepository) AddAssignments(assignments *models.Assignments) (bool, error) {
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
)`

	result, err := assign.db.NamedExec(query, assignments)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (assign *assignmentsRepository) GetAssignments() ([]models.Assignments, error) {
	const query = `
SELECT shift, number, receipt, amount
FROM green_seeds.assignments`

	var assignments []models.Assignments
	if err := assign.db.Select(&assignments, query); err != nil {
		return nil, err
	}

	return assignments, nil
}

func (assign *assignmentsRepository) UpdateAssignments(assignments *models.Assignments) (bool, error) {
	const query = `
UPDATE green_seeds.assignments
SET shift = :shift,
	number = :number,
	receipt = :receipt,
	amount = :amount
`

	result, err := assign.db.NamedExec(query, assignments)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (assign *assignmentsRepository) DeleteAssignments(id int64) (bool, error) {
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
