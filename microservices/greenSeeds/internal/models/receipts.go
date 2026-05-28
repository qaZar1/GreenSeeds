package models

import "time"

type Recipes struct {
	Recipe      *int64     `json:"recipe" db:"recipe"`
	Seed        string     `json:"seed" db:"seed"`
	SeedRu      string     `json:"seed_ru" db:"seed_ru"`
	Gcode       string     `json:"gcode" db:"gcode"`
	Updated     *time.Time `json:"updated" db:"updated"`
	Description string     `json:"description" db:"description"`
} // @name recipes
