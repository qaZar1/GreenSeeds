package service

import (
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (s *Service) AddReport(report models.Reports) (models.Reports, error) {
	if err := s.validate.Struct(report); err != nil {
		return models.Reports{}, err
	}

	return s.repo.RepRepo.AddReports(report)
}

func (s *Service) GetReports() ([]models.Reports, error) {
	return s.repo.RepRepo.GetReports()
}

func (s *Service) GetReportsByReport(idStr string) (models.Reports, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.Reports{}, err
	}

	return s.repo.RepRepo.GetReportsById(id)
}
