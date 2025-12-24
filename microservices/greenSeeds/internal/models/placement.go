package models

type Placement struct {
	Bunker int     `json:"bunker" db:"bunker"`
	Seed   string  `json:"seed" db:"seed"`
	SeedRu *string `json:"seed_ru" db:"seed_ru"`
	Amount uint64  `json:"amount" db:"amount"`
} // @name placement

type FillPlacement struct {
	Seed    string `json:"seed" db:"seed"`
	Percent int    `json:"percent" db:"percent"`
} // @name fillPlacement
