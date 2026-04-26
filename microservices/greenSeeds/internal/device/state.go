package device

type DeviceState string

const (
	StateUnknown      DeviceState = "UNKNOWN"
	StateManual       DeviceState = "MANUAL MODE"
	StateStandby      DeviceState = "STAND BY"
	StateReady        DeviceState = "READY"
	StateBusy         DeviceState = "BUSY"
	StateWait         DeviceState = "WAIT"
	StateReturn       DeviceState = "RETURN"
	StateError        DeviceState = "ERROR"
	StateDisconnected DeviceState = "DISCONNECTED"
)

type ManagerState string

const (
	ManagerStateConnecting   ManagerState = "CONNECTING"
	ManagerStateConnected    ManagerState = "CONNECTED"
	ManagerStateDisconnected ManagerState = "DISCONNECTED"
)

type ConnEventType int

const (
	ConnEventDisconnected ConnEventType = iota
	ConnEventConnected
)

type ConnEvent struct {
	Type ConnEventType
	Err  error
}
