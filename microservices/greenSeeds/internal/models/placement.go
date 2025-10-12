package models

type Placement struct {
	Bunker int    `json:"bunker" db:"bunker"`
	Seed   string `json:"seed" db:"seed"`
} // @name Placement
