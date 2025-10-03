package models

import "time"

type Receipts struct {
	Receipt     int       `json:"receipt" db:"receipt"`
	Seed        string    `json:"seed" db:"seed"`
	Gcode       string    `json:"gcode" db:"gcode"`
	Updated     time.Time `json:"updated" db:"updated"`
	Description string    `json:"description" db:"description"`
}
