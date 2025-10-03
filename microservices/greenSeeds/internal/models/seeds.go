package models

type Seeds struct {
	Seed         string `json:"seed" db:"seed"`
	MinDensity   int    `json:"min_density" db:"min_density"`
	MaxDensity   int    `json:"max_density" db:"max_density"`
	TankCapacity int    `json:"tank_capacity" db:"tank_capacity"`
	Latency      int    `json:"latency" db:"latency"`
}
