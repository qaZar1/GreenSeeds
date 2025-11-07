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
	CheckActiveTasks(username string) ([]models.ActiveTask, error)
	GetTaskById(id int) (models.Task, error)
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

func (assign *assignmentsRepository) CheckActiveTasks(username string) ([]models.ActiveTask, error) {
	const query = `
SELECT
	a.id,
    a.shift,
    a.number,
    a.receipt,
    s.dt,
    a.amount,
    COALESCE(SUM(CASE WHEN r.success THEN 1 ELSE 0 END), 0) AS done_turns,
    se.seed
FROM green_seeds.assignments a
JOIN green_seeds.shifts s ON s.shift = a.shift
LEFT JOIN green_seeds.reports r 
  ON r.shift = a.shift
  AND r.number = a.number 
  AND r.receipt = a.receipt
left join green_seeds.receipts r2 
  on r2.receipt = a.receipt
left join green_seeds.seeds se
  on se.seed = r2.seed 
WHERE s.username = $1 and DATE(s.dt) = CURRENT_DATE
GROUP BY a.id, a.shift, a.number, a.receipt, s.dt, a.amount, se.seed
HAVING COALESCE(SUM(CASE WHEN r.success THEN 1 ELSE 0 END), 0) < a.amount;
`

	var activeTasks []models.ActiveTask
	if err := assign.db.Select(&activeTasks, query, username); err != nil {
		return nil, err
	}

	return activeTasks, nil
}

func (assign *assignmentsRepository) GetTaskById(id int) (models.Task, error) {
	const query = `
select
	a.id,
    a.shift,
    a.number,
    r.seed,
	p.bunker,
	r.gcode,
    a.amount AS required_amount,
    COALESCE(SUM(CASE WHEN rp.success THEN 1 ELSE 0 END), 0) AS completed_amount
FROM green_seeds.assignments a
JOIN green_seeds.receipts r ON r.receipt = a.receipt
LEFT JOIN green_seeds.reports rp
    ON rp.shift = a.shift
    AND rp.number = a.number
LEFT JOIN green_seeds.placement p
	ON p.seed = r.seed
WHERE a.id = $1
GROUP BY a.id, a.shift, a.number, r.seed, p.bunker, r.gcode, a.amount;`

	var task models.Task
	if err := assign.db.Get(&task, query, id); err != nil {
		return models.Task{}, err
	}

	return task, nil
}
