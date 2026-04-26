package repository

import (
	"errors"
	"time"

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
	insReports(tx *sqlx.Tx, reports []models.Reports) (bool, error)
	delReports(tx *sqlx.Tx, assignments models.Assignments) (int64, error)
	SyncReports(oldAssignments, assignments models.Assignments, reports []models.Reports) (models.Assignments, error)
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
WHERE green_seeds.shifts.dt >= (CURRENT_DATE AT TIME ZONE 'UTC+5') AND a.deleted_at IS NULL
ORDER BY a.shift, a.number;`

	var assignments []models.Assignments
	if err := assign.db.Select(&assignments, query); err != nil {
		return nil, err
	}

	return assignments, nil
}

func (assign *assignmentsRepository) SyncReports(
	assignments models.Assignments,
	oldAssignments models.Assignments,
	reports []models.Reports,
) (models.Assignments, error) {
	var finished, pending []models.Reports

	for _, report := range reports {
		if report.Success || report.Error != nil {
			finished = append(finished, report)
		} else {
			pending = append(pending, report)
		}
	}

	tx, err := assign.db.Beginx()
	if err != nil {
		return models.Assignments{}, err
	}
	defer tx.Rollback()

	if assignments.Amount > oldAssignments.Amount {
		res := assignments.Amount - oldAssignments.Amount
		startTurn := oldAssignments.Amount + 1
		var reports []models.Reports
		for i := range res {
			reports = append(reports, models.Reports{
				Shift:   assignments.Shift,
				Number:  assignments.Number,
				Receipt: assignments.Receipt,
				Turn:    startTurn + i,
			})
		}

		if _, err := assign.insReports(tx, reports); err != nil {
			return models.Assignments{}, err
		}
	} else if assignments.Amount < oldAssignments.Amount {
		amount, err := assign.delReports(tx, assignments)
		if err != nil {
			return models.Assignments{}, err
		}

		assignments.Amount = oldAssignments.Amount - int(amount)
	}

	updated, err := assign.UpdateAssignments(tx, assignments)
	if err != nil {
		return models.Assignments{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Assignments{}, err
	}

	return updated, nil
}

func (assign *assignmentsRepository) UpdateAssignments(tx *sqlx.Tx, assignments models.Assignments) (models.Assignments, error) {
	const query = `
UPDATE green_seeds.assignments
SET shift = :shift,
	number = :number,
	receipt = :receipt,
	amount = :amount
WHERE id = :id AND deleted_at IS NULL
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
UPDATE green_seeds.assignments
SET deleted_at = $1
WHERE id = $2 AND deleted_at IS NULL;
`

	result, err := assign.db.Exec(query, time.Now(), id)
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
WHERE id = $1 AND deleted_at IS NULL;`

	var assignments models.Assignments
	if err := assign.db.Get(&assignments, query, id); err != nil {
		return models.Assignments{}, err
	}

	return assignments, nil
}

func (assign *assignmentsRepository) CheckActiveTasks(userId string) ([]models.ActiveTask, error) {
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
WHERE s.user_id = $1 and DATE(s.dt) = CURRENT_DATE AND a.deleted_at IS NULL
GROUP BY a.id, a.shift, a.number, a.receipt, s.dt, a.amount, se.seed, se.seed_ru 
HAVING COALESCE(SUM(CASE WHEN r.success THEN 1 ELSE 0 END), 0) < a.amount;
`

	var activeTasks []models.ActiveTask
	if err := assign.db.Select(&activeTasks, query, userId); err != nil {
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
    (
        SELECT p.bunker
        FROM green_seeds.placement p
        WHERE p.seed = r.seed
          AND p.bunker IS NOT NULL
        ORDER BY p.bunker
        LIMIT 1
    ) AS bunker,
    r.gcode,
    r.receipt,
    a.amount AS required_amount
FROM green_seeds.assignments a
JOIN green_seeds.receipts r
    ON r.receipt = a.receipt
LEFT JOIN green_seeds.seeds se
    ON r.seed = se.seed
WHERE a.id = $1 AND a.deleted_at IS NULL;`

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

func (assign *assignmentsRepository) insReports(tx *sqlx.Tx, reports []models.Reports) (bool, error) {
	const insertReportsQuery = `
INSERT INTO green_seeds.reports (
	shift,
	number,
	receipt,
	turn
)
VALUES (:shift, :number, :receipt, :turn)`

	result, err := tx.NamedExec(insertReportsQuery, reports)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == int64(len(reports)), nil
}

func (assign *assignmentsRepository) delReports(tx *sqlx.Tx, assignments models.Assignments) (int64, error) {
	const deleteReportsQuery = `
DELETE FROM green_seeds.reports
WHERE shift = :shift AND 
	number = :number AND
	receipt = :receipt AND
	turn > :amount AND
	success IS NULL AND
	error IS NULL`

	result, err := tx.NamedExec(
		deleteReportsQuery,
		assignments,
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
