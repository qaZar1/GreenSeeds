package device

import (
	"fmt"
	"strings"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (c *DeviceClient) Ping(sessionId string) error {
	defer c.RefreshPolling()

	matchResults := func(data []byte) MatchResult {
		str := strings.TrimSpace(string(data))
		c.Manager.parseStatus(str)
		switch {
		case strings.Contains(str, "ACK BOOT"):
			return MatchResult{Matched: true, Done: true}
		}
		return MatchResult{}
	}

	ch, err := c.Manager.Do(c.ctx, []byte("BOOT"+delimeterStr), matchResults, 10, sessionId)
	if err != nil {
		response := models.WSResponse{
			Type:  BOOT,
			Message: err.Error(),
		}
		select {
		case c.RespCh <- response:
		default:
		}

		return err
	}

	for {
		select {
		case resp, ok := <-ch:
			if !ok {
				return fmt.Errorf("channel closed")
			}
			respStr := string(resp)
			response := models.WSResponse{
				Type:   BOOT,
				Status: respStr,
			}
			select {
			case c.RespCh <- response:
			default:
			}

			return nil
		case <-c.ctx.Done():
			return nil
		}
	}
}

func (c *DeviceClient) Begin(sessionId string, command string, iter int) error {
	matchResults := func(data []byte) MatchResult {
		str := strings.TrimSpace(string(data))
		c.Manager.parseStatus(str)

		switch {
		case strings.Contains(str, "ACK"):
			return MatchResult{Matched: true}
		case strings.Contains(str, "END"):
			return MatchResult{Matched: true, Done: true}
		}
		return MatchResult{}
	}

	// run command
	ch, err := c.Manager.Do(c.ctx, []byte(command+delimeterStr), matchResults, 10, sessionId)
	if err != nil {
		c.Manager.ReleaseSession(sessionId)
		c.RefreshPolling()
		return err
	}

	// waiting results and transport to ws
	for resp := range ch {
		respStr := string(resp)

		respArr := strings.Split(respStr, " ")
		progress := strings.Split(respArr[1], "/")
		
		response := models.WSResponse{
			Type:   "STATE",
			Iteration: iter,
			Data: map[string]string{
				"shift":progress[0],
				"number":progress[1],
				"turn":progress[2],
			},
		}
		if respArr[0] == "ACK" {
			response.Status = "BEGIN_ACK"
		} else {
			response.Status = "BEGIN_END"
		}

		select {
		case c.RespCh <- response:
		default:
		}
	}

	return nil
}

func (c *DeviceClient) Return(sessionId string) error {
	matchResults := func(data []byte) MatchResult {
		str := strings.TrimSpace(string(data))
		c.Manager.parseStatus(str)

		switch {
		case strings.Contains(str, "STAND BY"):
			return MatchResult{Matched: true, Done: true}
		}
		return MatchResult{}
	}

	ch, err := c.Manager.Do(c.ctx, []byte("RETURN"+delimeterStr), matchResults, 10, sessionId)
	if err != nil {
		return err
	}

	for resp := range ch {
		respStr := string(resp)
		response := models.WSResponse{
			Type:   RETURN,
			Status: respStr,
		}

		select {
		case c.RespCh <- response:
		default:
		}
	}

	return nil
}

func (c *DeviceClient) GetStatus() {
	state := c.Manager.GetState()

	stStr := fmt.Sprintf("%v", state)

	response := models.WSResponse{
		Type:   STATUS,
		Status: stStr,
	}

	select {
	case c.RespCh <- response:
	default:
	}
}

func (c *DeviceClient) SetStatusReady(sessionId string) {
	defer c.RefreshPolling()

	matchResults := func(data []byte) MatchResult {
		str := strings.TrimSpace(string(data))
		c.Manager.parseStatus(str)

		switch {
		case strings.Contains(str, ACK_SETSTATUS_READY):
			return MatchResult{Matched: true, Done: true}
		}
		return MatchResult{}
	}

	ch, err := c.Manager.Do(c.ctx, []byte(SETSTATUS_READY+delimeterStr), matchResults, 10, sessionId)
	if err != nil {
		response := models.WSResponse{
			Type:  SETSTATUS_READY,
			Message: err.Error(),
		}

		select {
		case c.RespCh <- response:
		default:
		}
	}

	for resp := range ch {
		c.Manager.SetState(StateReady)
		respStr := string(resp)
		response := models.WSResponse{
			Type:   SETSTATUS_READY,
			Status: respStr,
		}

		select {
		case c.RespCh <- response:
		default:
		}
	}
}

func (c *DeviceClient) Status(sessionId string) {
	defer c.RefreshPolling()

	matchResults := func(data []byte) MatchResult {
		str := strings.TrimSpace(string(data))
		c.Manager.parseStatus(str)
		
		return MatchResult{Matched: true, Done: true}
	}

	ch, err := c.Manager.Do(c.ctx, []byte(STATUS+delimeterStr), matchResults, 10, sessionId)
	if err != nil {
		response := models.WSResponse{
			Type:  STATUS,
			Status: "ERROR",
			Message: err.Error(),
		}

		select {
		case c.RespCh <- response:
		default:
		}
	}

	for resp := range ch {
		respStr := string(resp)
		response := models.WSResponse{
			Type:   STATUS,
			Status: "OK",
			Message: respStr,
		}

		select {
		case c.RespCh <- response:
		default:
		}
	}
}

func (c *DeviceClient) Stop(sessionId string) error {
	c.Manager.ReleaseSession(sessionId)
	c.RefreshPolling()
	return nil
}
