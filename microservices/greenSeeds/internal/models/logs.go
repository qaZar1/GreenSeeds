package models

import "time"

type Log struct {
	Id        int64     `db:"id" json:"id"`
	Dt        time.Time `db:"dt" json:"dt"`
	Lvl       string    `db:"lvl" json:"lvl"`
	RequestId *string   `db:"request_id" json:"request_id"`
	Msg       string    `db:"msg" json:"msg"`
	Call      *string   `db:"caller" json:"caller"`
	Username  *string   `db:"username" json:"username"`
	UserId    *int64    `db:"user_id" json:"user_id"`
} // @name log

type LogsParams struct {
	Search   string     `json:"search" db:"search"`
	Level    string     `json:"level" db:"level"`
	Limit    string     `json:"limit" db:"limit"`
	Offset   string     `json:"offset" db:"offset"`
	DateFrom *time.Time `json:"date_from" db:"date_from"`
	DateTo   *time.Time `json:"date_to" db:"date_to"`
}
