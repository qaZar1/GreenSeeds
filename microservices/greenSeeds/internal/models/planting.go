package models

import (
	"bytes"
)


type Planting struct {
	Active       bool
	Iteration    int
	MaxIter      int
	Bunker		 int

	PrevState    State
	CurrentState State
	Command      string
	Seed 		 string
	Error        error

	LastBuf    *bytes.Buffer
	LastRespAI *ResponseAPI

	Params Params

	Results *[]Reports
	ActiveBunker *SeedsWithBunker

	Stop bool
}
