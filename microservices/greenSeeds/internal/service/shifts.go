package service

import (
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (s *Service) AddShift(shifts models.Shifts) (models.Shifts, error) {
	if err := s.validate.Struct(shifts); err != nil {
		return models.Shifts{}, err
	}

	return s.repo.ShfRepo.AddShifts(shifts)
}

func (s *Service) GetShifts() ([]models.Shifts, error) {
	return s.repo.ShfRepo.GetShifts()
}

func (s *Service) GetShiftsByShift(shift string) (models.Shifts, error) {
	shiftInt, err := strconv.Atoi(shift)
	if err != nil {
		return models.Shifts{}, err
	}

	return s.repo.ShfRepo.GetShiftsByShift(shiftInt)
}

func (s *Service) UpdateShifts(shifts models.Shifts) (models.Shifts, error) {
	if err := s.validate.Struct(shifts); err != nil {
		return models.Shifts{}, err
	}

	return s.repo.ShfRepo.UpdateShifts(shifts)
}

func (s *Service) DeleteShifts(shift string) (bool, error) {
	shiftInt, err := strconv.Atoi(shift)
	if err != nil {
		return false, err
	}

	return s.repo.ShfRepo.DeleteShifts(shiftInt)
}

func (s *Service) GetShiftsWithoutUser() ([]models.Shifts, error) {
	return s.repo.ShfRepo.GetShiftsWithoutUser()
}
