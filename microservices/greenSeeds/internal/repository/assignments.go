package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IAssignmentsRepository interface {
	AddAssignments(assignments models.Assignments) (models.Assignments, error)
	GetAssignments() ([]models.Assignments, error)
	UpdateAssignments(tx *sqlx.Tx, assignments models.Assignments) (models.Assignments, error)
	DeleteAssignments(number int) (bool, error)
	GetAssignmentsByNumber(number int) (models.Assignments, error)
	CheckActiveTasks(username string) ([]models.ActiveTask, error)
	GetTaskById(id int) (models.Task, error)
	insReports(tx *sqlx.Tx, oldAssignments, assignments models.Assignments) (bool, error)
	delReports(tx *sqlx.Tx, oldAssignments, assignments models.Assignments) (bool, error)
	Transaction(oldAssignments, assignments models.Assignments) models.Assignments
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

	var reports []models.Reports
	for i := range inserted.Amount {
		reports = append(reports, models.Reports{
			Shift:   inserted.Shift,
			Number:  inserted.Number,
			Receipt: inserted.Receipt,
			Turn:    i + 1,
		})
	}

	const insertReportsQuery = `
INSERT INTO green_seeds.reports (
	shift,
	number,
	receipt,
	turn
)
VALUES (:shift, :number, :receipt, :turn)`
	res, err := assign.db.NamedExec(insertReportsQuery, reports)
	if err != nil {
		return models.Assignments{}, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Assignments{}, err
	}

	if rowsAffected != int64(len(reports)) {
		return models.Assignments{}, errors.New("some reports were not inserted")
	}

	return inserted, nil
}

func (assign *assignmentsRepository) GetAssignments() ([]models.Assignments, error) {
	const query = `
SELECT 
	id,
	a.shift as shift,
	number,
	a.receipt as receipt,
	r.description as description,
	amount
FROM green_seeds.assignments a
JOIN green_seeds.shifts ON a.shift = green_seeds.shifts.shift
LEFT JOIN green_seeds.receipts r
ON r.receipt = a.receipt
WHERE green_seeds.shifts.dt >= (CURRENT_DATE AT TIME ZONE 'UTC+5')
ORDER BY a.shift, a.number;`

	var assignments []models.Assignments
	if err := assign.db.Select(&assignments, query); err != nil {
		return nil, err
	}

	return assignments, nil
}

func (assign *assignmentsRepository) Transaction(
	assignments models.Assignments,
	oldAssignments models.Assignments,
) models.Assignments {
	tx, err := assign.db.Beginx()
	if err != nil {
		return models.Assignments{}
	}
	defer tx.Rollback()

	updated, err := assign.UpdateAssignments(tx, assignments)
	if err != nil {
		return models.Assignments{}
	}

	if oldAssignments.Amount > assignments.Amount {
		ok, err := assign.delReports(tx, oldAssignments, assignments)
		if err != nil {
			return models.Assignments{}
		}
		if !ok {
			return models.Assignments{}
		}
	} else if oldAssignments.Amount < assignments.Amount {
		ok, err := assign.insReports(tx, oldAssignments, assignments)
		if err != nil {
			return models.Assignments{}
		}
		if !ok {
			return models.Assignments{}
		}
	}

	if err := tx.Commit(); err != nil {
		return models.Assignments{}
	}

	return updated
}

func (assign *assignmentsRepository) UpdateAssignments(tx *sqlx.Tx, assignments models.Assignments) (models.Assignments, error) {
	const query = `
UPDATE green_seeds.assignments
SET shift = :shift,
	number = :number,
	receipt = :receipt,
	amount = :amount
WHERE id = :id
RETURNING id, shift, number, receipt, amount;
`

	rows, err := tx.NamedQuery(query, assignments)
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
    se.seed,
    se.seed_ru
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
GROUP BY a.id, a.shift, a.number, a.receipt, s.dt, a.amount, se.seed, se.seed_ru 
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
SELECT
	a.id,
    a.shift,
    a.number,
    r.seed,
    se.seed_ru,
	p.bunker,
	r.gcode,
	r.receipt,
    a.amount AS required_amount
FROM green_seeds.assignments a
JOIN green_seeds.receipts r
ON r.receipt = a.receipt
LEFT JOIN green_seeds.reports rp
    ON rp.shift = a.shift
    AND rp.number = a.number
LEFT JOIN green_seeds.placement p
	ON p.seed = r.seed
LEFT JOIN green_seeds.seeds se
	ON r.seed = se.seed 
WHERE a.id = $1
GROUP BY a.id, a.shift, a.number, r.seed, se.seed_ru, p.bunker, r.gcode, r.receipt, a.amount;`

	var task models.Task
	if err := assign.db.Get(&task, query, id); err != nil {
		return models.Task{}, err
	}

	const reportsQuery = `
SELECT
	r.shift,
	r.number,
	r.receipt,
	r.turn,
	COALESCE(r.success, FALSE) AS success,
	r.error,
	r.solution,
	r.mark
FROM green_seeds.reports r
WHERE r.shift = $1 AND r.number = $2 AND r.receipt = $3 AND (r.success = FALSE OR r.success IS NULL)
ORDER BY r.turn ASC`
	var reports []models.Reports
	if err := assign.db.Select(&reports, reportsQuery, task.Shift, task.Number, task.Receipt); err != nil {
		return models.Task{}, err
	}

	task.Reports = &reports

	return task, nil
}

func (assign *assignmentsRepository) insReports(tx *sqlx.Tx, oldAssignments, assignments models.Assignments) (bool, error) {
	const insertReportsQuery = `
INSERT INTO green_seeds.reports (
	shift,
	number,
	receipt,
	turn
)
VALUES (:shift, :number, :receipt, :turn)`

	var reports []models.Reports
	for i := oldAssignments.Amount; i < assignments.Amount; i++ {
		reports = append(reports, models.Reports{
			Shift:   assignments.Shift,
			Number:  assignments.Number,
			Receipt: assignments.Receipt,
			Turn:    i + 1,
		})
	}

	result, err := tx.NamedExec(insertReportsQuery, reports)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == int64(assignments.Amount-oldAssignments.Amount), nil
}

func (assign *assignmentsRepository) delReports(tx *sqlx.Tx, oldAssignments, assignments models.Assignments) (bool, error) {
	const deleteReportsQuery = `
DELETE FROM green_seeds.reports
WHERE shift = :shift AND number = :number AND receipt = :receipt AND turn > :amount`

	result, err := tx.NamedExec(
		deleteReportsQuery,
		assignments,
	)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == int64(oldAssignments.Amount-assignments.Amount), nil
}
