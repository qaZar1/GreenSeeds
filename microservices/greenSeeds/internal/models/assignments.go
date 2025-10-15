package models

type Assignments struct {
	Id      *int64 `json:"id" db:"id"`
	Shift   int64  `json:"shift" db:"shift"`
	Number  int    `json:"number" db:"number"`
	Receipt int64  `json:"receipt" db:"receipt"`
	Amount  int    `json:"amount" db:"amount"`
}
