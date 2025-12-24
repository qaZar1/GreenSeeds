package models

import (
	"time"
)

type Assignments struct {
	Id          *int64 `json:"id" db:"id"`
	Shift       int64  `json:"shift" db:"shift"`
	Number      int    `json:"number" db:"number"`
	Receipt     int64  `json:"receipt" db:"receipt"`
	ReceiptDesc string `json:"description" db:"description"`
	Amount      int    `json:"amount" db:"amount"`
} // @name assignment

type ActiveTask struct {
	Id        int64     `json:"id" db:"id"`
	Shift     int64     `json:"shift" db:"shift"`
	Number    int       `json:"number" db:"number"`
	Receipt   int64     `json:"receipt" db:"receipt"`
	Dt        time.Time `json:"dt" db:"dt"`
	Amount    int       `json:"amount" db:"amount"`
	DoneTurns int       `json:"done_turns" db:"done_turns"`
	Seed      string    `json:"seed" db:"seed"`
	SeedRu    string    `json:"seed_ru" db:"seed_ru"`
} // @name active_task

type Task struct {
	Id             int64      `json:"id" db:"id"`
	Shift          int64      `json:"shift" db:"shift"`
	Number         int        `json:"number" db:"number"`
	Seed           string     `json:"seed" db:"seed"`
	SeedRu         string     `json:"seed_ru" db:"seed_ru"`
	Bunker         int64      `json:"bunker" db:"bunker"`
	Gcode          string     `json:"gcode" db:"gcode"`
	Receipt        int64      `json:"receipt" db:"receipt"`
	RequiredAmount int        `json:"required_amount" db:"required_amount"`
	Reports        *[]Reports `json:"reports" db:"reports"`
} // @name task
