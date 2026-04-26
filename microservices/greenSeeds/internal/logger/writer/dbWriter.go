package writer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/rs/zerolog"
)

type DbWriter struct {
	DB *sqlx.DB
}

func NewDbWriter(db *sqlx.DB) *DbWriter {
	return &DbWriter{DB: db}
}

func (w *DbWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (w *DbWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level < zerolog.InfoLevel {
		return len(p), nil
	}

	var log models.Log
	if err := json.Unmarshal(p, &log); err != nil {
		fmt.Fprintf(os.Stderr, "json unmarshal: %v\n", err)
		return len(p), err
	}

	if log.Call == nil {
		nilStr := ""
		log.Call = &nilStr
	}

	if log.Username == nil {
		nilStr := ""
		log.Username = &nilStr
	}

	w.DB.Get(&log.UserId, "SELECT id FROM green_seeds.users WHERE username = :username", log.Username)

	log.Lvl = strings.ToUpper(log.Lvl)

	const query = `
INSERT INTO green_seeds.logs (dt, lvl, request_id, msg, caller, user_id)
VALUES (:dt, :lvl, :request_id, :msg, :caller, :user_id);`

	result, err := w.DB.NamedExec(query, log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "db insert error: %v\n", err)
		return len(p), err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return len(p), err
	}

	if rowsAffected == 0 {
		fmt.Fprintf(os.Stderr, "db insert error: %v\n", err)
		return len(p), err
	}

	return len(p), nil
}
