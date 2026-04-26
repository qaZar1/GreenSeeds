package models

type WSMessageType string

const (
	TypeBoot      WSMessageType = "BOOT"
	TypeBegin     WSMessageType = "BEGIN"
	TypeReturn    WSMessageType = "RETURN"
	TypeStatus    WSMessageType = "STATUS"
	TypeError     WSMessageType = "ERROR"
	TypeSetStatus WSMessageType = "SET STATUS READY"
	TypeAuth      WSMessageType = "AUTH"

	TypeRetry WSMessageType = "RETRY"
	TypeSkip  WSMessageType = "SKIP"
	TypeAbort WSMessageType = "ABORT"

	TypeState   = "STATE"   // шаги процесса
    TypeDevice  = "DEVICE"  // статус устройства
    TypeAction  = "ACTION"  // запрос действия
    TypeEnd     = "END"
)

type State int

type Action string