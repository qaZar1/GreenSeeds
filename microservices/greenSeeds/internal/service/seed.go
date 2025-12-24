package service

import (
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (s *Service) AddSeed(seed models.Seeds) (models.Seeds, error) {
	if err := s.validate.Struct(seed); err != nil {
		return models.Seeds{}, err
	}

	return s.repo.SeedRepo.AddSeeds(seed)
}

func (s *Service) GetSeeds() ([]models.Seeds, error) {
	return s.repo.SeedRepo.GetSeeds()
}

func (s *Service) GetSeedBySeed(seed string) (models.Seeds, error) {
	return s.repo.SeedRepo.GetSeedsBySeed(seed)
}

func (s *Service) GetSeedWithBunkers(seed string) ([]models.SeedsWithBunker, error) {
	return s.repo.SeedRepo.GetSeedsWithBunkers(seed)
}

func (s *Service) UpdateSeed(seed models.Seeds) (models.Seeds, error) {
	if err := s.validate.Struct(seed); err != nil {
		return models.Seeds{}, err
	}

	return s.repo.SeedRepo.UpdateSeeds(seed)
}

func (s *Service) DeleteSeed(seed string) (bool, error) {
	return s.repo.SeedRepo.DeleteSeeds(seed)
}
