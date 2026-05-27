package ws

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/device"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func waitForDeviceReady(manager *device.Manager, c *Client, timeout time.Duration, iter *models.Iteration) error {
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

		case msg, ok := <-c.Control:
			if !ok {
				return errors.New("receive channel closed")
			}

			if msg.Type == "STOP" {
				if iter != nil {
					iter.Finished = true
				}
				if c != nil {
					c.planting.Stop = true
					safeSend(c.Send, okResponse("STOP", "Stopped by user"))
				}
				return errors.New("Stopped")
			}

		case <-ticker.C:
			if manager.GetStatus() != device.ManagerStateConnected {
				return errors.New("device disconnected")
			}

			if manager.GetState() == device.StateReady {
				return nil // Успех
			}
		}
	}
}

func RunIteration(s *Server, c *Client, iter *models.Iteration) {
	iter.CurrentState = StateWaitingReady

	for {

		if !isDeviceAlive(s.dClient.Manager) {
			Err(iter, errors.New("device disconnected"))
			break
		}

		if checkStop(c) {
			Err(iter, errors.New("stopped by user"))
		}

		if s.dClient.Manager.GetState() == device.StateDisconnected {
			c.planting.Stop = true
			Err(iter, errors.New("device state disconnected"))
			break
		}

		switch iter.CurrentState {

		case StateWaitingReady:
			state(
				s,
				StatusWaitReady,
				MessageWaitingReady,
				iter.Turn,
			)

			if err := waitForDeviceReady(
				s.dClient.Manager,
				c,
				30*time.Second,
				iter,
			); err != nil {
				Err(iter, err)
				break
			}

			iter.CurrentState = StateBegin

		case StateBegin:
			state(
				s,
				StatusBegin,
				MessageBegin,
				iter.Turn,
			)

			if err := Begin(s, c, iter); err != nil {
				Err(iter, err)
				break
			}

			state(
				s,
				StatusBegin,
				MessageSeedPlanted,
				iter.Turn,
			)

			iter.CurrentState = StatePhoto

		case StatePhoto:
			state(
				s,
				StatusPhoto,
				MessagePhoto,
				iter.Turn,
			)

			if err := Photo(s, c, iter); err != nil {
				Err(iter, err)
				break
			}

			state(
				s,
				StatusPhoto,
				MessagePhotoDone,
				iter.Turn,
			)

			iter.CurrentState = StateControl

		case StateControl:
			state(
				s,
				StatusControl,
				MessageControlStart,
				iter.Turn,
			)

			ok, err := Control(s, c, iter)
			if err != nil {
				Err(iter, err)
				break
			}

			if !ok {
				iter.CurrentState = StateBegin
				continue
			}

			state(
				s,
				StatusControl,
				MessageControlEnd,
				iter.Turn,
			)

			iter.CurrentState = StateProcess

		case StateProcess:
			state(
				s,
				StatusReturn,
				MessageWaitingReturn,
				iter.Turn,
			)

			if err := s.dClient.Return(c.SessionId); err != nil {
				Err(iter, err)
				break
			}

			state(
				s,
				StatusReturn,
				MessageReturned,
				iter.Turn,
			)

			iter.CurrentState = StateDone

		case StateDone:
			iter.Success = true
			iter.Finished = true
		}

		if iter.Finished {
			errsStr := ErrorsToString(iter.Err)

			finishIteration(
				s,
				c,
				iter,
				iter.Success,
				errsStr,
				"",
				"",
			)

			if iter.Success {
				state(s, StatusDone, MessageDone, iter.Turn)
			} else {
				state(s, StatusError, errsStr, iter.Turn)
			}

			c.planting.Stop = true

			return
		}
	}
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

func ArrayToString(strs []string) string {
	builder := strings.Builder{}
	for _, str := range strs {
		builder.WriteString(str + ";\n")
	}

	return builder.String()
}

func Begin(s *Server, c *Client, iter *models.Iteration) error {
	bunkers, err := s.repo.SeedRepo.GetBestBunker(iter.Seed)
	if err != nil {
		err = fmt.Errorf("bunkers with seed %s are empty", iter.Seed)

		safeSend(c.Send, errResponse(models.TypeBegin, err))

		c.planting.Stop = true

		return err
	}

	iter.Bunker = bunkers.Bunker

	command := buildGcode(c, iter)

	if err := s.dClient.Begin(c.SessionId, command, iter.Turn); err != nil {
		return err
	}

	if err := s.repo.PlcRepo.DecrementSeed(iter.Bunker); err != nil {
		safeSend(c.Send, errResponse(models.TypeBegin, errors.New("bunker is empty")))
		return err
	}

	return nil
}

func Photo(s *Server, c *Client, iter *models.Iteration) error {
	buf, err := s.camera.GetBytesFromPhoto("./Proj_img/Amarant/photo_2024-06-18_15-22-18.jpg")
	if err != nil {
		return err
	}

	if buf == nil || buf.Len() == 0 {
		return emptyPhotoError()
	}

	iter.LastBuf = buf

	return nil
}

func Control(s *Server, c *Client, iter *models.Iteration) (bool, error) {
	response, err := s.api.CheckAI(iter.Seed, *iter.LastBuf)
	if err != nil {
		return false, err
	}

	if response.PercentOfMatch < s.config.MinMatchPercent {
		return false, errors.New(
			fmt.Sprintf("percent of match is less than %d", int(s.config.MinMatchPercent*100)),
		)
	}

	if response.Seed != iter.Seed {
		return false, errors.New(
			fmt.Sprintf("seed mismatch. expected %s, but got %s", iter.Seed, response.Seed),
		)
	}

	count := s.opencv.Counter(iter.LastBuf.Bytes())
	if count == 0 {
		return false, errors.New("counting failed")
	}

	seed, err := s.repo.SeedRepo.GetSeedsBySeed(iter.Seed)
	if err != nil {
		return false, err
	}

	if count < seed.MinDensity || count > seed.MaxDensity {
		return false, nil
	}

	return true, nil
}
