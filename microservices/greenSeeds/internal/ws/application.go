package ws

import (
	"errors"
	"strings"
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func handleBoot(s *Server, client *Client, req models.WSRequest) {
	if err := s.dClient.Ping(client.SessionId); err != nil {
		return
	}
}

func handleStatus(s *Server, c *Client, req models.WSRequest) {
	if !s.dClient.Manager.TryAcquireSession(c.SessionId) {
		c.Send <- errResponse("START", errors.New("device busy"))
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
		safeSend(c.Send, errResponse("START", errors.New("params required")))
		return
	}

	// Пытаемся захватить устройство
	if !s.dClient.Manager.TryAcquireSession(c.SessionId) {
		safeSend(c.Send, errResponse("START", errors.New("device busy")))
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
		safeSend(c.Send, errResponse("START", errors.New("Cant start")))
		return
	}

	recipe, err := s.repo.RptRepo.GetRecipesByRecipe(req.Params.Recipe)
	if err != nil {
		safeSend(c.Send, errResponse("START", errors.New("Cant start")))
		return
	}

	// Инициализация
	c.planting = models.Planting{
		Active:    true,
		Iteration: 0,
		MaxIter:   len(allReports) - 1,
	}

	for c.planting.Iteration < c.planting.MaxIter {
		if c.planting.Stop {
			break
		}

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

		c.planting.Iteration++
		time.Sleep(10 * time.Second)
	}

	// Всё закончили
	safeSend(c.Send, okResponse("END", "Посадка завершена"))
}

func handleAuth(s *Server, client *Client, req models.WSRequest) {
	if req.Token == nil {
		safeSend(client.Send, errResponse(req.Type, errors.New("Token is nil")))
		return
	}

	token := *req.Token
	if !strings.HasPrefix(token, "Bearer ") {
		safeSend(client.Send, errResponse(req.Type, errors.New("Invalid token format")))
		return
	}

	claims, err := s.infra.GetTokenClaims(token[7:])
	if err != nil {
		safeSend(client.Send, errResponse(req.Type, errors.New("Failed to get token claims")))
		return
	}

	if claims == nil {
		safeSend(client.Send, errResponse(req.Type, errors.New("Invalid token")))
		return
	}

	if err := validateJWT(*claims, s.repo); err != nil {
		safeSend(client.Send, errResponse(req.Type, errors.New("Invalid token")))
		return
	}

	client.IsAuth = true

	safeSend(client.Send, okResponse(req.Type, "OK"))
}
