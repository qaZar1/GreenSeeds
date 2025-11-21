package models

type ResponseAPI struct {
	Seed           string  `json:"seed"`
	PercentOfMatch float64 `json:"percent_of_match"`
}
