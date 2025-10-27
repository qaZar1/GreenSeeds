package models

import (
	"time"
)

type Assignments struct {
	Id      *int64 `json:"id" db:"id"`
	Shift   int64  `json:"shift" db:"shift"`
	Number  int    `json:"number" db:"number"`
	Receipt int64  `json:"receipt" db:"receipt"`
	Amount  int    `json:"amount" db:"amount"`
}

type ActiveTask struct {
	Id        int64     `json:"id" db:"id"`
	Shift     int64     `json:"shift" db:"shift"`
	Number    int       `json:"number" db:"number"`
	Receipt   int64     `json:"receipt" db:"receipt"`
	Dt        time.Time `json:"dt" db:"dt"`
	Amount    int       `json:"amount" db:"amount"`
	DoneTurns int       `json:"done_turns" db:"done_turns"`
	Seed      string    `json:"seed" db:"seed"`
}

type Task struct {
	Id              int64  `json:"id" db:"id"`
	Shift           int64  `json:"shift" db:"shift"`
	Number          int    `json:"number" db:"number"`
	Seed            string `json:"seed" db:"seed"`
	RequiredAmount  int    `json:"required_amount" db:"required_amount"`
	CompletedAmount int    `json:"completed_amount" db:"completed_amount"`
}
