package models

import "time"

type Shifts struct {
	Shift    *int64    `json:"shift" db:"shift"`
	Dt       time.Time `json:"dt" db:"dt"`
	Username *string   `json:"username" db:"username"`
} //@name Shift
