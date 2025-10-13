package service

import (
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (s *Service) AddBunker(bunker models.Bunkers) (models.Bunkers, error) {
	if err := s.validate.Struct(bunker); err != nil {
		return models.Bunkers{}, err
	}

	return s.repo.BunkRepo.AddBunkers(bunker)
}

func (s *Service) GetBunkers() ([]models.Bunkers, error) {
	return s.repo.BunkRepo.GetBunkers()
}

func (s *Service) GetBunkersById(bunkerId string) (models.Bunkers, error) {
	bunkerIdInt, err := strconv.Atoi(bunkerId)
	if err != nil {
		return models.Bunkers{}, err
	}

	return s.repo.BunkRepo.GetBunkersById(bunkerIdInt)
}

func (s *Service) UpdateBunker(bunker models.Bunkers) (models.Bunkers, error) {
	if err := s.validate.Struct(bunker); err != nil {
		return models.Bunkers{}, err
	}

	return s.repo.BunkRepo.UpdateBunkers(bunker)
}

func (s *Service) DeleteBunker(bunkerId string) (bool, error) {
	bunkerIdInt, err := strconv.Atoi(bunkerId)
	if err != nil {
		return false, err
	}

	return s.repo.BunkRepo.DeleteBunkers(bunkerIdInt)
}
