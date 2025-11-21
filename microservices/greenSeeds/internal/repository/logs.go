package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type ILogsRepository interface {
	GetLogs(params models.LogsParams) ([]models.Log, error)
}

type logsRepository struct {
	db *sqlx.DB
}

func NewLogsRepository(db *sqlx.DB) *logsRepository {
	return &logsRepository{
		db: db,
	}
}

func (lg *logsRepository) GetLogs(params models.LogsParams) ([]models.Log, error) {
	query := `
    SELECT 
        id,
		dt,
		lvl,
		request_id,
		msg,
		caller,
		username
    FROM green_seeds.logs
    WHERE 1=1`

	if params.Level != "" && params.Level != "ALL" {
		query += fmt.Sprintf(" AND lvl = :level")
	}

	if params.Search != "" {
		query += fmt.Sprintf(` AND (
            msg ILIKE :search OR 
            request_id ILIKE :search OR 
            caller ILIKE :search OR 
            username ILIKE :search
        )`)
	}

	if params.DateFrom != nil {
		query += fmt.Sprintf(" AND dt >= :date_from")
	}

	if params.DateTo != nil {
		query += fmt.Sprintf(" AND dt <= :date_to")
	}

	query += " ORDER BY dt DESC"

	query += fmt.Sprintf(" LIMIT :limit OFFSET :offset")

	rows, err := lg.db.NamedQuery(query, params)
	if err != nil {
		return nil, fmt.Errorf("error running named query: %w", err)
	}
	defer rows.Close()

	var logs []models.Log
	for rows.Next() {
		var log models.Log
		if err := rows.StructScan(&log); err != nil {
			return nil, fmt.Errorf("error scanning log row: %w", err)
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return logs, nil
}
