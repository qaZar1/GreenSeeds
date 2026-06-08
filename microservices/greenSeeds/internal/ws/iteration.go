package ws

import (
	"fmt"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func RunIteration(s *Server, c *Client, iter *models.Iteration) {
	iter.CurrentState = StateWaitingReady
	iter.IsFailed = false
	iter.CountRetry = 0

	for {
		if iter.CurrentState != StateProcess {
			if !isDeviceAlive(s.dClient.Manager) {
				AddError(s.log, iter, fmt.Errorf("Устройство отключено"), "")

				iter.IsFailed = true
				iter.CurrentState = StateProcess
				continue
			}
		}

		switch iter.CurrentState {
		case StateWaitingReady:
			stateWaitingReady(s, c, iter)

		case StateBegin:
			stateBegin(s, c, iter)

		case StatePhoto:
			statePhoto(s, c, iter)
			continue

		case StateControl:
			stateControl(s, c, iter)
			iter.CurrentState = StateProcess
			continue

		case StateProcess:
			stateProcess(s, c, iter)

		case StateDone:
			stateDone(iter)
		}

		if iter.Finished {
			err := ""
			if iter.Err != nil {
				err = iter.Err.Error()
			}
			finishIteration(
				s,
				c,
				iter,
				err,
				"",
				"",
			)

			if iter.Success == models.ReportStatusSuccess {
				EmitDone(c, MessageDone, iter, c.planting.Required)
			} else {
				EmitError(
					c,
					MessageError,
					iter,
					iter.Err,
					iter.ErrStage,
					c.planting.Required,
				)
			}

			return
		}
	}
}
