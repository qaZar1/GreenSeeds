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

	StatusWaitReady = "WAIT_READY"
	StatusBegin     = "BEGIN"
	StatusPhoto     = "PHOTO"
	StatusControl   = "CONTROL"
	StatusReturn    = "RETURN"
	StatusError     = "ERROR"
	StatusDone      = "DONE"

	MessageWaitingReady  = "Ожидание готовности устройства"
	MessageBegin         = "Начало процесса посева"
	MessageSeedPlanted   = "Посадка выполнена, переход к фотографированию"
	MessagePhoto         = "Фотографирование"
	MessagePhotoDone     = "Фотографирование завершено"
	MessageControlStart  = "Начало проведения контроля качества"
	MessageControlEnd    = "Контроль качества завершен"
	MessageControlFailed = "Контроль качества не пройден"
	MessageWaitingReturn = "Возврат устройства"
	MessageReturned      = "Устройство возвращено"
	MessageError         = "Ошибка"
	MessageDone          = "Процесс завершен"
)
