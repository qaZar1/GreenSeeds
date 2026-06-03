package ws

type EventType string

const (
	EventState EventType = "STATE"
	EventError EventType = "ERROR"
	EventDone  EventType = "DONE"
	EventRetry EventType = "RETRY"
	EventStop  EventType = "STOP"
)

type Step string

const (
	StepWaitReady Step = "WAIT_READY"
	StepBegin     Step = "BEGIN"
	StepPhoto     Step = "PHOTO"
	StepControl   Step = "CONTROL"
	StepProcess   Step = "PROCESS"
	StepReturn    Step = "RETURN"
	StepDone      Step = "DONE"
)
