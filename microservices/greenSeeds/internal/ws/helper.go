package ws

import "github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"

func okResponse(t models.WSMessageType, msg string) models.WSResponse {
	return models.WSResponse{
		Type:    t,
		Status:  "OK",
		Message: msg,
	}
}

func errResponse(t models.WSMessageType, err error) models.WSResponse {
	return models.WSResponse{
		Type:    t,
		Status:  "ERROR",
		Message: err.Error(),
	}
}

func state(
	s *Server,
	status string,
	message string,
	iter int,
) {
	safeSend(s.Send, models.WSResponse{
		Type:      "STATE",
		Status:    status,
		Message:   message,
		Iteration: iter,
	})
}

func safeSend(ch chan models.WSResponse, msg models.WSResponse) {
	select {
	case ch <- msg:
	default:
	}
}
