package ws

import (
	"log"
	"strings"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func OkResponse(
	t models.WSMessageType,
	msg string,
) models.WSResponse {
	return models.WSResponse{
		Type:    t,
		Status:  "OK",
		Message: msg,
	}
}

func ErrResponse(
	t models.WSMessageType,
	err error,
) models.WSResponse {
	return models.WSResponse{
		Event:   string(t),
		Status:  "ERROR",
		Message: err.Error(),

		Error: &models.WSError{
			Code:    mapErrorCode(err),
			Message: err.Error(),
		},
	}
}

func Emit(
	c *Client,
	event string,
	msg string,
	iter *models.Iteration,
) {
	c.Send <- models.WSResponse{
		Event:     event,
		Message:   msg,
		Iteration: iter.Turn,
	}
}

func EmitState(
	c *Client,
	step string,
	msg string,
	iter *models.Iteration,
	total int,
) {
	SafeSend(c.Send, models.WSResponse{
		Event:     "STATE",
		Message:   msg,
		Iteration: iter.Turn,
		Step:      step,

		Progress: &models.Progress{
			Current: iter.Turn,
			Total:   total,
			Percent: calcPercent(iter.Turn, total),
		},
	})
}

func EmitError(
	c *Client,
	step string,
	iter *models.Iteration,
	err error,
	stage string,
	total int,
) {
	resp := models.WSResponse{
		Event:   "ERROR",
		Message: err.Error(),
		Step:    step,

		Error: &models.WSError{
			Code:    mapErrorCode(err),
			Stage:   stage,
			Message: err.Error(),
		},
	}

	if iter != nil {
		resp.Iteration = iter.Turn

		resp.Progress = &models.Progress{
			Current: iter.Turn,
			Total:   total,
			Percent: calcPercent(iter.Turn, total),
		}
	}

	c.Send <- resp
}

func EmitDone(
	c *Client,
	msg string,
	iter *models.Iteration,
	total int,
) {
	c.Send <- models.WSResponse{
		Event:     "DONE",
		Message:   msg,
		Iteration: iter.Turn,

		Progress: &models.Progress{
			Current: total,
			Total:   total,
			Percent: 100,
		},
	}
}

func mapErrorCode(err error) string {
	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "timeout"):
		return "TIMEOUT"

	case strings.Contains(msg, "device disconnected"):
		return "DEVICE_DISCONNECTED"

	case strings.Contains(msg, "seed mismatch"):
		return "SEED_MISMATCH"

	case strings.Contains(msg, "density"):
		return "DENSITY_INVALID"

	default:
		return "UNKNOWN"
	}
}

func SafeSend(ch chan models.WSResponse, msg models.WSResponse) {
	select {
	case ch <- msg:
	default:
		log.Println("Failed to send message to client")
	}
}
