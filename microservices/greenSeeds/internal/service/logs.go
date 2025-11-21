package service

import (
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (s *Service) GetLogs(params models.LogsParams) ([]models.Log, error) {
	return s.repo.LogsRepo.GetLogs(params)
}
