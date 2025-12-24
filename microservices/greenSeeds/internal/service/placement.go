package service

import (
	"errors"
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (s *Service) AddPlacement(placement models.Placement) (models.Placement, error) {
	if err := s.validate.Struct(placement); err != nil {
		return models.Placement{}, err
	}

	seed, err := s.GetSeedBySeed(placement.Seed)
	if err != nil {
		return models.Placement{}, err
	}

	if placement.Amount > uint64(seed.TankCapacity) {
		return models.Placement{}, errors.New("amount is greater than tank capacity")
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

	seed, err := s.GetSeedBySeed(placement.Seed)
	if err != nil {
		return models.Placement{}, err
	}

	if placement.Amount > uint64(seed.TankCapacity) {
		return models.Placement{}, errors.New("amount is greater than tank capacity")
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

func (s *Service) FillPlacment(fillPlacement models.FillPlacement) (models.Placement, error) {
	seedWithBunkers, err := s.GetSeedWithBunkers(fillPlacement.Seed)
	if err != nil {
		return models.Placement{}, err
	}

	if len(seedWithBunkers) == 0 {
		return models.Placement{}, errors.New("Seed not found")
	}

	min := seedWithBunkers[0]
	for _, bunker := range seedWithBunkers {
		if bunker.Amount < min.Amount {
			min = bunker
		}
	}

	amount := min.TankCapacity * fillPlacement.Percent / 100
	placement := models.Placement{
		Bunker: min.Bunker,
		Seed:   min.Seed,
		Amount: uint64(amount),
	}

	return s.repo.PlcRepo.UpdatePlacement(placement)
}
