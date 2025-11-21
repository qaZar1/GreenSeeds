package api

import (
	"bytes"
	"math/rand/v2"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

var seeds = map[int]string{
	1: "Amaranth",
	2: "Basil",
	3: "Watercress",
	4: "Melissa",
	5: "Borago",
}

func (a *API) CheckAI(seed string, buf bytes.Buffer) (models.ResponseAPI, error) {
	_ = seed
	_ = buf
	seedNumber := rand.IntN(4)
	percent := rand.Float64()
	return models.ResponseAPI{
		Seed:           seeds[seedNumber+1],
		PercentOfMatch: percent,
	}, nil
}
