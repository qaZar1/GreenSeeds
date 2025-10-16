package service

import (
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

// func (s *Service) AddAssignment(assignment models.Assignments) (models.Assignments, error) {
// 	if err := s.validate.Struct(assignment); err != nil {
// 		return models.Assignments{}, err
// 	}

// 	return s.repo.AsnRepo.AddAssignments(assignment)
// }

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

// func (s *Service) UpdateAssignment(assignment models.Assignments) (models.Assignments, error) {
// 	if err := s.validate.Struct(assignment); err != nil {
// 		return models.Assignments{}, err
// 	}

// 	return s.repo.AsnRepo.UpdateAssignments(assignment)
// }

// func (s *Service) DeleteAssignments(idStr string) (bool, error) {
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return false, err
// 	}

// 	return s.repo.AsnRepo.DeleteAssignments(id)
// }
