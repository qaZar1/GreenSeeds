package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	_ "modernc.org/sqlite"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLiteClient(cfg models.Config) *SQLite {
	db, err := sql.Open("sqlite", "./db/calibration.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	sqlite := &SQLite{db: db}

	go sqlite.startCleaner()

	return sqlite
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) startCleaner() {
	ticker := time.NewTicker(1 * time.Hour)

	for range ticker.C {
		s.Cleaner()
	}
}

func (s *SQLite) Cleaner() {
	calibrations, err := s.GetOldCalibration()
	if err != nil {
		fmt.Printf("Invalid get old calibration: %s", err)
		return
	}

	for _, calibration := range calibrations {
		if calibration.FirstPhotoPath == nil {
			continue
		}

		dir := (*calibration.FirstPhotoPath)[:strings.LastIndex(*calibration.FirstPhotoPath, "/")]
		os.RemoveAll(dir)
	}

	if err := s.DeleteOldRows(); err != nil {
		log.Fatal(err)
	}

	return
}

func (s *SQLite) AddCalibration(calibration models.Calibration) (bool, error) {
	const query = `
INSERT INTO calibration (
	session_id,
	created_at
)
VALUES (
	$1,
	$2
);`

	res, err := s.db.Exec(query, calibration.SessionId, calibration.CreatedAt)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (s *SQLite) UpdateCalibration(calibration models.Calibration) (bool, error) {
	const query = `
UPDATE calibration
SET first_photo_path = COALESCE($1, first_photo_path),
	second_photo_path = COALESCE($2, second_photo_path),
	dx = COALESCE($3, dx),
	dy = COALESCE($4, dy),
	cir = COALESCE($5, cir),
	d_per_step = COALESCE($6, d_per_step)
WHERE session_id = $7;`

	res, err := s.db.Exec(
		query,
		calibration.FirstPhotoPath,
		calibration.SecondPhotoPath,
		calibration.Dx, calibration.Dy,
		calibration.Cir,
		calibration.DPerStep,
		calibration.SessionId,
	)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (s *SQLite) GetCalibration(sessionId string) (models.Calibration, error) {
	const query = `
SELECT
	session_id,
	first_photo_path,
	second_photo_path,
	dx,
	dy,
	cir,
	d_per_step,
	created_at
FROM calibration
WHERE session_id = $1;`

	var calibration models.Calibration
	if err := s.db.QueryRow(
		query,
		sessionId,
	).Scan(
		&calibration.SessionId,
		&calibration.FirstPhotoPath,
		&calibration.SecondPhotoPath,
		&calibration.Dx,
		&calibration.Dy,
		&calibration.Cir,
		&calibration.DPerStep,
		&calibration.CreatedAt,
	); err != nil {
		return models.Calibration{}, err
	}

	return calibration, nil
}

func (s *SQLite) GetOldCalibration() ([]models.Calibration, error) {
	const query = `
SELECT
	session_id,
	first_photo_path,
	second_photo_path,
	dx,
	dy,
	cir,
	d_per_step,
	created_at
FROM calibration
WHERE created_at < $1;
`
	rows, err := s.db.Query(query, time.Now().AddDate(0, 0, -1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Calibration

	for rows.Next() {
		var c models.Calibration
		if err := rows.Scan(
			&c.SessionId,
			&c.FirstPhotoPath,
			&c.SecondPhotoPath,
			&c.Dx,
			&c.Dy,
			&c.Cir,
			&c.DPerStep,
			&c.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, nil
}

func (s *SQLite) DeleteOldRows() error {
	const query = `
DELETE FROM calibration
WHERE created_at < $1;`

	_, err := s.db.Exec(query, time.Now().AddDate(0, 0, -1))
	if err != nil {
		return err
	}

	return nil
}
