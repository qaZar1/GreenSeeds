package ws

import (
	"errors"
	"strings"
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func handleBoot(s *Server, client *Client, req models.WSRequest) {
	if err := s.dClient.Boot(client.SessionId, true); err != nil {
		return
	}
}

func handleStatus(s *Server, c *Client, req models.WSRequest) {
	if !s.dClient.Manager.TryAcquireSession(c.SessionId) {
		c.Send <- ErrResponse("START", errors.New("device busy"))
		s.dClient.GetStatus()

		return
	}

	defer s.dClient.Manager.ReleaseSession(c.SessionId)

	// Останавливаем фоновый пинг на время сессии
	s.dClient.PausePolling()
	defer s.dClient.RefreshPolling()

	s.dClient.Status(c.SessionId)
}

func handleSetStatusReady(s *Server, client *Client, req models.WSRequest) {
	s.dClient.SetStatusReady(client.SessionId)
}

func handlePlanting(s *Server, c *Client, req models.WSRequest) {
	// Проверяем параметры
	if req.Params == nil {
		SafeSend(c.Send, ErrResponse("START", errors.New("params required")))
		return
	}

	// Пытаемся захватить устройство
	if !s.dClient.Manager.TryAcquireSession(c.SessionId) {
		SafeSend(c.Send, ErrResponse("START", errors.New("device busy")))
		return
	}
	// Гарантированно освобождаем при выходе из функции
	defer s.dClient.Manager.ReleaseSession(c.SessionId)

	// Останавливаем фоновый пинг на время сессии
	s.dClient.PausePolling()
	defer s.dClient.RefreshPolling()

	allReports, err := s.repo.RepRepo.GetNotSuccessfulAssignments(
		req.Params.Shift,
		req.Params.Number,
		req.Params.Recipe,
	)
	if err != nil {
		SafeSend(c.Send, ErrResponse("START", errors.New("Cant start")))
		return
	}

	recipe, err := s.repo.RptRepo.GetRecipesByRecipe(req.Params.Recipe)
	if err != nil {
		SafeSend(c.Send, ErrResponse("START", errors.New("Cant start")))
		return
	}

	// Инициализация
	c.planting = models.Planting{
		Active:    true,
		Iteration: 0,
		MaxIter:   len(allReports) - 1,
		StopChan:  make(chan struct{}),
		Required:  req.Params.RequiredAmount,
	}

	for c.planting.Iteration < c.planting.MaxIter {
		report := allReports[c.planting.Iteration]
		iter := models.Iteration{
			Seed:      recipe.Seed,
			Gcode:     recipe.Gcode,
			Shift:     req.Params.Shift,
			Number:    req.Params.Number,
			Turn:      report.Turn,
			Required:  req.Params.RequiredAmount,
			Recipe:    int(*recipe.Recipe),
			ExtraMode: req.Params.ExtraMode,
			Report:    report,
		}

		RunIteration(s, c, &iter)

		if !iter.Success {
			break
		}

		c.planting.Iteration++

		select {
		case <-c.planting.StopChan:
			Emit(c, "END", "Автоматическая посадка прервана", &models.Iteration{})
			return
		case <-time.After(10 * time.Second):
		}
	}

	// Всё закончили
	Emit(c, "END", "Автоматическая посадка завершена", &models.Iteration{})
}

func handleStop(s *Server, c *Client, req models.WSRequest) {
	select {
	case <-c.planting.StopChan:
		// уже закрыт
	default:
		close(c.planting.StopChan)
	}

	SafeSend(c.Send, OkResponse("STOP", "Stopping after current iteration"))
}

func handleAuth(s *Server, client *Client, req models.WSRequest) {
	if req.Token == nil {
		SafeSend(client.Send, ErrResponse(req.Type, errors.New("Token is nil")))
		return
	}

	token := *req.Token
	if !strings.HasPrefix(token, "Bearer ") {
		SafeSend(client.Send, ErrResponse(req.Type, errors.New("Invalid token format")))
		return
	}

	claims, err := s.infra.GetTokenClaims(token[7:])
	if err != nil {
		SafeSend(client.Send, ErrResponse(req.Type, errors.New("Failed to get token claims")))
		return
	}

	if claims == nil {
		SafeSend(client.Send, ErrResponse(req.Type, errors.New("Invalid token")))
		return
	}

	if err := validateJWT(*claims, s.repo); err != nil {
		SafeSend(client.Send, ErrResponse(req.Type, errors.New("Invalid token")))
		return
	}

	client.IsAuth = true

	SafeSend(client.Send, OkResponse(req.Type, "OK"))
}
