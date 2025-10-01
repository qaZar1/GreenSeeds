package models

import "time"

type Shifts struct {
	Shift    int       `json:"shift" db:"shift"`
	Dt       time.Time `json:"dt" db:"dt"`
	Username string    `json:"username" db:"username"`
}
