package service

import (
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (s *Service) AddPlacement(placement models.Placement) (models.Placement, error) {
	if err := s.validate.Struct(placement); err != nil {
		return models.Placement{}, err
	}

	return s.repo.PlcRepo.AddPlacement(placement)
}

func (s *Service) GetPlacements() ([]models.Placement, error) {
	return s.repo.PlcRepo.GetPlacement()
}

func (s *Service) GetPlacementByBunker(bunkerId string) (models.Placement, error) {
	bunkerIdInt, err := strconv.Atoi(bunkerId)
	if err != nil {
		return models.Placement{}, err
	}

	return s.repo.PlcRepo.GetPlacementByBunker(bunkerIdInt)
}

func (s *Service) UpdatePlacement(placement models.Placement) (models.Placement, error) {
	if err := s.validate.Struct(placement); err != nil {
		return models.Placement{}, err
	}

	return s.repo.PlcRepo.UpdatePlacement(placement)
}

func (s *Service) DeletePlacement(bunkerId string) (bool, error) {
	bunkerIdInt, err := strconv.Atoi(bunkerId)
	if err != nil {
		return false, err
	}

	return s.repo.PlcRepo.DeletePlacement(bunkerIdInt)
}
