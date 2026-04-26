package models

import (
	"bytes"
)

type Iteration struct {
	PrevState    State
	CurrentState State
	NextState State
	
	Seed string
	Gcode string
	Bunker int
	Shift int
	Number int
	Turn int
	Required int
	Receipt int

	ExtraMode bool

	Err []error
	LastBuf *bytes.Buffer

	Report Reports

	Solutions []string

	Success bool
	Finished bool
}