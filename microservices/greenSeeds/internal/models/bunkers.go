package models

type Bunkers struct {
	Bunker   int `json:"bunker" db:"bunker"`
	Distance int `json:"distance" db:"distance"`
}
