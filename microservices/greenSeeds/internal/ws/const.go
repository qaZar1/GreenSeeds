package ws

import "github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"

const (
	StateWaitingReady models.State = iota
	StateBegin
	StatePhoto
	StateControl
	StateWaitingReturn
	StateProcess
	StateError
	StateDone
)

const (
	ActionSkip  models.Action = "SKIP"
	ActionRetry models.Action = "RETRY"
	ActionAbort models.Action = "ABORT"
)