package service

import (
	"errors"
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (s *Service) AddAssignment(assignment models.Assignments) (models.Assignments, error) {
	if err := s.validate.Struct(assignment); err != nil {
		return models.Assignments{}, err
	}

	return s.repo.AsnRepo.AddAssignments(assignment)
}

func (s *Service) GetAssignments() ([]models.Assignments, error) {
	return s.repo.AsnRepo.GetAssignments()
}

func (s *Service) GetAssignmentsByAssignment(idStr string) (models.Assignments, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.Assignments{}, err
	}

	return s.repo.AsnRepo.GetAssignmentsByNumber(id)
}

func (s *Service) UpdateAssignment(assignment models.Assignments) (models.Assignments, error) {
	if err := s.validate.Struct(assignment); err != nil {
		return models.Assignments{}, err
	}

	oldAssignment, err := s.repo.AsnRepo.GetAssignmentsByNumber(int(*assignment.Id))
	if err != nil {
		return models.Assignments{}, err
	}

	updated := s.repo.AsnRepo.Transaction(assignment, oldAssignment)
	if updated == (models.Assignments{}) {
		return models.Assignments{}, errors.New("transaction failed")
	}

	return updated, nil
}

func (s *Service) DeleteAssignments(idStr string) (bool, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return false, err
	}

	return s.repo.AsnRepo.DeleteAssignments(id)
}

func (s *Service) CheckActiveTasks(username string) ([]models.ActiveTask, error) {
	return s.repo.AsnRepo.CheckActiveTasks(username)
}

func (s *Service) GetTaskById(idStr string) (models.Task, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.Task{}, err
	}

	return s.repo.AsnRepo.GetTaskById(id)
}
