package ws

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/device"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

// В том же файле или в helpers
func waitForDeviceReady(manager *device.Manager, c *Client, timeout time.Duration) error {
	ticker := time.NewTicker(500 * time.Millisecond) // опрос раз в 0.5с
	defer ticker.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return errors.New("timeout waiting for READY status")

		case <-c.done:
			return errors.New("client disconnected")

		case _, ok := <-c.Control:
			if !ok {
				return errors.New("receive channel closed")
			}

		case <-ticker.C:
			if manager.GetStatus() != device.ManagerStateConnected {
				return errors.New("device disconnected")
			}

			if manager.GetStatus() == device.READY {
				return nil // Успех
			}
		}
	}
}

func RunIteration(s *Server, c *Client, iter *models.Iteration) {
	iter.PrevState = StateWaitingReady
	iter.CurrentState = StateWaitingReady
	iter.NextState = StateBegin

	for{
		if !isDeviceAlive(s.dClient.Manager){
			iter.Success = false
			iter.Finished = true

			iter.Err = append(iter.Err, errors.New("device disconnected"))
			return
		}
		if checkStop(c) {
			iter.Success = false
			iter.Finished = true
		}

		if s.dClient.Manager.GetState() == device.StateDisconnected{
			c.planting.Stop = true
			iter.Finished = true
		}
		
		switch iter.CurrentState {
		case StateWaitingReady:
			state(s, "WAIT_READY", iter.Turn)

			if err := waitForDeviceReady(s.dClient.Manager, c, 30*time.Second); err != nil {
				continue
			}

			iter.CurrentState = StateBegin

		case StateBegin:
			Begin(s, c, iter)
		case StatePhoto:
			Photo(s, c, iter)

			state(s, "PHOTO_DONE", iter.Turn)
		case StateControl:
			Control(s, c, iter)

			state(s, "AI_OK", iter.Turn)
		case StateProcess:
			if checkStop(c) {
				c.planting.Stop = true
				iter.CurrentState = StateDone
				return
			}

			if err := s.dClient.Return(c.SessionId); err != nil{
				iter.PrevState = StateControl
				iter.CurrentState = StateError
				iter.Err = append(iter.Err, err)
				continue
			}

			state(s, "RETURN_DONE", iter.Turn)

			iter.CurrentState = StateDone

		case StateError:
			solution := Err(s, c, iter)
			iter.Solutions = append(iter.Solutions, solution)
		case StateDone:
			iter.Success = true
			iter.Finished = true
		}

		if iter.Finished{
			errsStr := ErrorsToString(iter.Err)
			solutionsStr := ArrayToString(iter.Solutions)
			
			finishIteration(
				s,
				c,
				iter,
				iter.Success,
				errsStr,
				solutionsStr,
				"",
			)

			state(s, "DONE", iter.Turn)
			return
		}
	}
}

func isDeviceAlive(m *device.Manager) bool {
    return m.GetStatus() == device.ManagerStateConnected
}

func buildGcode(c *Client, iter *models.Iteration) string {
	const query = `
BEGIN %d/%d/%d
BUNKER %d
\x02%s\x03`

	data := fmt.Sprintf(query,
		iter.Shift,
		iter.Number,
		iter.Turn,
		iter.Bunker,
		iter.Gcode)

	if iter.ExtraMode {
		sorrel := buildSorrel(iter)
		data = fmt.Sprintf("%s%s", data, sorrel)
	}

	return data
}

func buildSorrel(iter *models.Iteration) string {
	const query = `
\x07Sorrel\x0A%d/%d\x0A%d/%d\x0A\x0D`

	data := fmt.Sprintf(query,
		iter.Shift,
		iter.Number,
		iter.Turn,
		iter.Required)
	return data
}

func buildReport(
	c *Client,
	dt time.Time,
	success bool,
	err string,
	solution string,
	mark string,
	iter *models.Iteration,
) models.Reports {
	return models.Reports{
		Shift:    int64(iter.Shift),
		Number:   iter.Number,
		Receipt:  int64(iter.Receipt),
		Turn:     iter.Turn,
		Dt:       &dt,
		Success:  success,
		Error:    &err,
		Solution: &solution,
		Mark:     &mark,
	}
}

func finishIteration(
	s *Server,
	c *Client,
	iter *models.Iteration,
	success bool,
	errStr string,
	solution string,
	mark string,
) {
	now := time.Now()

	report := buildReport(
		c,
		now,
		success,
		errStr,
		solution,
		mark,
		iter,
	)

	s.repo.RepRepo.UpdateReports(report)
}

func SetActiveBunker(bunkersBySeed []models.SeedsWithBunker) *models.SeedsWithBunker {
	if len(bunkersBySeed) == 0 {
		return nil
	}

	var selected *models.SeedsWithBunker

	for i := range bunkersBySeed {
		b := &bunkersBySeed[i]

		// пропускаем пустые
		if b.Amount <= 0 {
			continue
		}

		if selected == nil || b.Amount > selected.Amount {
			selected = b
		}
	}

	return selected
}

func ErrorsToString(errs []error) string {
	builder := strings.Builder{}
	for _, err := range errs {
		builder.WriteString(err.Error() + ";\n")
	}

	return builder.String()
}

func ArrayToString(strs []string) string {
	builder := strings.Builder{}
	for _, str := range strs {
		builder.WriteString(str + ";\n")
	}

	return builder.String()
}

func checkStop(c *Client) bool {
    select {
    case <-c.Control:
        c.planting.Stop = true
        c.Send <- okResponse("STOP", "Stopped by user")
        return true
    default:
        return false
    }
}

func waitAction(c *Client) (models.Action, string, bool) {
    select {
    case msg := <-c.Actions:
        switch msg.Type {
        case "RETRY":
            return ActionRetry, "RETRY", true
        case "SKIP":
            return ActionSkip, "SKIP", true
        case "ABORT":
            return ActionAbort, "ABORT", false
        }

		return ActionAbort, "", false
    case <-time.After(60 * time.Second):
        return ActionAbort, "", false
    }

}

func Begin(s *Server, c *Client, iter *models.Iteration) {
	bunkers, err := s.repo.SeedRepo.GetBestBunker(iter.Seed)
	if err != nil{
		c.Send <- errResponse(models.TypeBegin, errors.New("bunker not validate"))
		return
	}

	iter.Bunker = bunkers.Bunker
	
	command := buildGcode(c, iter)

	if err := s.dClient.Begin(c.SessionId, command, iter.Turn); err != nil {
		iter.PrevState = StateBegin
		iter.CurrentState = StateError
		iter.Err = append(iter.Err, err)
		return
	}
	
	if err := s.repo.PlcRepo.DecrementSeed(iter.Bunker); err != nil {
		c.Send <- errResponse("BEGIN", errors.New("bunker is empty"))

		iter.PrevState = StateBegin
		iter.CurrentState = StateError
		iter.Err = append(iter.Err, err)
		return
	}

	iter.CurrentState = StatePhoto
}

func Photo(s *Server, c *Client, iter *models.Iteration) {
	// buf, err = s.camera.TakePhoto()
	buf, err := s.camera.GetBytesFromPhoto("./Proj_img/Amarant/photo_2024-06-18_15-22-18.jpg")
	if err != nil || buf == nil || buf.Len() == 0 {
		iter.PrevState = StatePhoto
		iter.CurrentState = StateError
		iter.Err = append(iter.Err, err)
		return
	}

	iter.LastBuf = buf
	iter.CurrentState = StateControl
}

func Control(s *Server, c *Client, iter *models.Iteration) {
	_, err := s.api.CheckAI(iter.Seed, *iter.LastBuf)
	if err != nil {
		iter.PrevState = StateControl
		iter.CurrentState = StateError
		iter.Err = append(iter.Err, err)
		return
	}

	// if resp.Seed != c.planting.Params.Seed {
	// 	err = errors.New("wrong seed")
	// 	iter.PrevState = StateControl
	// 	iter.CurrentState = StateError
	// 	iter.Err = append(iter.Err, err)
	// 	return
	// }

	iter.CurrentState = StateProcess
}

func Err(s *Server, c *Client, iter *models.Iteration) string {
	if len(iter.Err) == 0{
		return ""
	}

	lastErr := iter.Err[len(iter.Err)-1]
	s.Send <- errResponseWithActions(
		models.TypeBegin,
		lastErr,
		[]string{"RETRY", "ABORT"},
		iter.Turn,
	)

	action, lastSolution, ok := waitAction(c)

	if !ok {
		iter.Success = false
		iter.Finished = true
		iter.Solutions = append(iter.Solutions, lastSolution)
		return ""
	}

	switch action {
	case ActionRetry:
		iter.CurrentState = iter.PrevState
		return lastSolution

	case ActionSkip:
		iter.CurrentState = StateProcess
		return lastSolution

	case ActionAbort:
		iter.Success = false
		iter.Finished = true
		iter.Solutions = append(iter.Solutions, lastSolution)
		return ""
	}

	return lastSolution
}