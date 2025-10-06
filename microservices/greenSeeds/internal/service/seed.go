package service

import (
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (s *Service) AddSeed(seed models.Seeds) (bool, error) {
	if err := s.validate.Struct(seed); err != nil {
		return false, err
	}

	return s.repo.SeedRepo.AddSeeds(seed)
}

func (s *Service) GetSeeds() ([]models.Seeds, error) {
	return s.repo.SeedRepo.GetSeeds()
}

func (s *Service) GetSeedById(seed string) (models.Seeds, error) {
	return s.repo.SeedRepo.GetSeedsBySeed(seed)
}

func (s *Service) UpdateSeed(seed models.Seeds) (bool, error) {
	if err := s.validate.Struct(seed); err != nil {
		return false, err
	}

	return s.repo.SeedRepo.UpdateSeeds(seed)
}

func (s *Service) DeleteSeed(seed string) (bool, error) {
	return s.repo.SeedRepo.DeleteSeeds(seed)
}
