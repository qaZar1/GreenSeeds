package models

type Placement struct {
	Bunker int     `json:"bunker" db:"bunker"`
	Seed   string  `json:"seed" db:"seed"`
	SeedRu *string `json:"seed_ru" db:"seed_ru"`
} // @name Placement
