package models

import "time"

type Reports struct {
	Id          *int64     `json:"id" db:"id"`
	Shift       int64      `json:"shift" db:"shift"`
	Number      int        `json:"number" db:"number"`
	Receipt     int64      `json:"receipt" db:"receipt"`
	Turn        int        `json:"turn" db:"turn"`
	Dt          *time.Time `json:"dt" db:"dt"`
	Success     bool       `json:"success" db:"success"`
	Error       *string    `json:"error" db:"error"`
	Solution    *string    `json:"solution" db:"solution"`
	Mark        *string    `json:"mark" db:"mark"`
	Responsible string     `json:"responsible" db:"full_name"`
}
