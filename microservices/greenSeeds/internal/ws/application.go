package ws

import (
	"errors"
	"strings"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func handleBoot(s *Server, client *Client, req models.WSRequest) {
	if err := s.dClient.Ping(client.SessionId); err != nil {
		return
	}
}

func handleStatus(s *Server, client *Client, req models.WSRequest) {
	s.dClient.GetStatus()
}

func handleSetStatusReady(s *Server, client *Client, req models.WSRequest) {
	s.dClient.SetStatusReady(client.SessionId)
}

func handlePlanting(s *Server, c *Client, req models.WSRequest) {
	// 1. Проверяем параметры
	if req.Params == nil {
		c.Send <- errResponse("START", errors.New("params required"))
		return
	}

	// 2. Пытаемся захватить устройство
	if !s.dClient.Manager.TryAcquireSession(c.SessionId) {
		c.Send <- errResponse("START", errors.New("device busy"))
		return
	}
	// Гарантированно освобождаем при выходе из функции
	defer s.dClient.Manager.ReleaseSession(c.SessionId)

	// 3. Останавливаем фоновый пинг на время сессии
	s.dClient.PausePolling()
	defer s.dClient.RefreshPolling()

	//
	allReports, err := s.repo.RepRepo.GetNotSuccessfulAssignments(
		req.Params.Shift,
		req.Params.Number,
		req.Params.Receipt,
	)
	if err != nil {
		c.Send <- errResponse("START", errors.New("Cant start"))
		return
	}

	receipt, err := s.repo.RptRepo.GetReceiptsByReceipt(req.Params.Receipt)
	if err != nil {
		c.Send <- errResponse("START", errors.New("Cant start"))
		return
	}

	// 4. Инициализация
	c.planting = models.Planting{
		Active:       true,
		Iteration:    0,
		MaxIter:      len(allReports),
	}
	
	for c.planting.Iteration < c.planting.MaxIter {
		report := allReports[c.planting.Iteration]
		iter := models.Iteration{
			Seed:receipt.Seed,
			Gcode: receipt.Gcode,
			Shift: req.Params.Shift,
			Number: req.Params.Number,
			Turn: report.Turn,
			Required: req.Params.RequiredAmount,
			Receipt: int(*receipt.Receipt),
			ExtraMode: req.Params.ExtraMode,
			Report: report,
		}

		RunIteration(s, c, &iter)

		if c.planting.Stop{
			break
		}

		c.planting.Iteration++ 
	}

	// Всё закончили
	c.Send <- okResponse("END", "Stopped")
}

func handleAuth(s *Server, client *Client, req models.WSRequest) {
	if req.Token == nil {
		client.Send <- errResponse(req.Type, errors.New("Token is nil"))
		return
	}

	token := *req.Token
	if !strings.HasPrefix(token, "Bearer ") {
		client.Send <- errResponse(req.Type, errors.New("Invalid token format"))
		return
	}

	claims, err := s.infra.GetTokenClaims(token[7:])
	if err != nil {
		client.Send <- errResponse(req.Type, errors.New("Failed to get token claims"))
		return
	}

	if claims == nil {
		client.Send <- errResponse(req.Type, errors.New("Invalid token"))
		return
	}

	if err := validateJWT(*claims, s.repo); err != nil {
		client.Send <- errResponse(req.Type, errors.New("Invalid token"))
		return
	}

	client.IsAuth = true

	client.Send <- okResponse(req.Type, "OK")
}
