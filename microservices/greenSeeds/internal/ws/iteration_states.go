package ws

import (
	"fmt"
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func stateWaitingReady(s *Server, c *Client, iter *models.Iteration) {
	EmitState(c, StatusWaitReady, MessageWaitingReady, iter, c.planting.Required)

	if err := waitForDeviceReady(
		s.dClient.Manager,
		c,
		30*time.Second,
		iter,
	); err != nil {

		AddError(s.log, iter, fmt.Errorf("Устройство не получило команду"), string(StepWaitReady))

		iter.IsFailed = true
		iter.CurrentState = StateProcess

		return
	}

	iter.CurrentState = StateBegin
}

func stateBegin(s *Server, c *Client, iter *models.Iteration) {
	iter.NeedDosow = false
	EmitState(c, StatusBegin, MessageBegin, iter, c.planting.Required)

	if err := Begin(s, c, iter); err != nil {

		AddError(s.log, iter, fmt.Errorf("Ошибка начала процесса посева"), string(StepBegin))

		iter.IsFailed = true
		iter.CurrentState = StateProcess

		return
	}

	EmitState(c, StatusBegin, MessageSeedPlanted, iter, c.planting.Required)

	iter.CurrentState = StatePhoto
}

func statePhoto(s *Server, c *Client, iter *models.Iteration) {
	EmitState(c, StatusPhoto, MessagePhoto, iter, c.planting.Required)

	if err := Photo(s, c, iter); err != nil {
		AddError(s.log, iter, fmt.Errorf("Ошибка фотографирования"), string(StepPhoto))

		iter.IsFailed = true
		iter.CurrentState = StateProcess

		return
	}

	EmitState(c, StatusPhoto, MessagePhotoDone, iter, c.planting.Required)

	iter.CurrentState = StateControl
}

func stateControl(s *Server, c *Client, iter *models.Iteration) {
	EmitState(c, StatusControl, MessageControlStart, iter, c.planting.Required)

	ok, err := Control(s, c, iter)
	if err != nil {
		AddError(s.log, iter, fmt.Errorf("Ошибка контроля"), string(StepControl))
		iter.IsFailed = true
		return
	}

	if !ok {
		if iter.CountRetry >= 3 {
			AddError(s.log, iter, fmt.Errorf("Слишком много попыток контроля"), string(StepControl))

			iter.IsFailed = true

			return
		}

		iter.CountRetry++
		iter.NeedDosow = true
		EditGcode(iter)

		EmitState(c, StatusControl, "Требуется досеивание", iter, c.planting.Required)

		return
	}

	EmitState(c, StatusControl, MessageControlEnd, iter, c.planting.Required)
}

func stateProcess(s *Server, c *Client, iter *models.Iteration) {
	EmitState(c, StatusReturn, MessageWaitingReturn, iter, c.planting.Required)

	if err := s.dClient.Return(c.SessionId); err != nil {

		AddError(s.log, iter, fmt.Errorf("Ошибка возврата устройства"), string(StepReturn))

		iter.Finished = true

		return
	}

	EmitState(c, StatusReturn, MessageReturned, iter, c.planting.Required)

	if isStopped(c) {
		AddError(s.log, iter, fmt.Errorf("Процесс остановлен пользователем"), string(StepReturn))

		iter.Finished = true

		return
	}

	if iter.IsFailed {
		iter.Finished = true
		return
	}

	if iter.NeedDosow {
		iter.CurrentState = StateBegin
		return
	}

	iter.CurrentState = StateDone
}

func stateDone(iter *models.Iteration) {
	iter.Finished = true
	iter.Success = models.ReportStatusSuccess
}
