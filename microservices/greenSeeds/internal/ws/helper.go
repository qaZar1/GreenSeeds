package ws

import "github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"

func okResponse(t models.WSMessageType, msg string) models.WSResponse {
	return models.WSResponse{
		Type:   t,
		Status: "OK",
		Message: msg,
	}
}

func errResponse(t models.WSMessageType, err error) models.WSResponse {
	return models.WSResponse{
		Type:  t,
		Status: "ERROR",
		Message: err.Error(),
	}
}

func errResponseWithActions(t models.WSMessageType, err error, actions []string, iter int) models.WSResponse {
	return models.WSResponse{
		Type:   "ACTION",
		Status: "ERROR",
		Message:  err.Error(),
		Iteration: iter,
		Actions: &actions,
	}
}

func state(s *Server, status string, iter int) {
    s.Send <- models.WSResponse{
        Type:      "STATE",
        Status:    status,
        Iteration: iter,
    }
}
