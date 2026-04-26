package application

import (
	"errors"
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (app *App) AddPlacement(placement models.Placement) (models.Placement, error) {
	if err := app.validate.Struct(placement); err != nil {
		return models.Placement{}, err
	}

	seed, err := app.GetSeedBySeed(placement.Seed)
	if err != nil {
		return models.Placement{}, err
	}

	if placement.Amount > uint64(seed.TankCapacity) {
		return models.Placement{}, errors.New("amount is greater than tank capacity")
	}

	return app.repo.PlcRepo.AddPlacement(placement)
}

func (app *App) GetPlacements() ([]models.Placement, error) {
	return app.repo.PlcRepo.GetPlacement()
}

func (app *App) GetPlacementByBunker(bunkerId string) (models.Placement, error) {
	bunkerIdInt, err := strconv.Atoi(bunkerId)
	if err != nil {
		return models.Placement{}, err
	}

	return app.repo.PlcRepo.GetPlacementByBunker(bunkerIdInt)
}

func (app *App) UpdatePlacement(placement models.Placement) (models.Placement, error) {
	if err := app.validate.Struct(placement); err != nil {
		return models.Placement{}, err
	}

	seed, err := app.GetSeedBySeed(placement.Seed)
	if err != nil {
		return models.Placement{}, err
	}

	if placement.Amount > uint64(seed.TankCapacity) {
		return models.Placement{}, errors.New("amount is greater than tank capacity")
	}

	return app.repo.PlcRepo.UpdatePlacement(placement)
}

func (app *App) DeletePlacement(bunkerId string) (bool, error) {
	bunkerIdInt, err := strconv.Atoi(bunkerId)
	if err != nil {
		return false, err
	}

	return app.repo.PlcRepo.DeletePlacement(bunkerIdInt)
}

func (app *App) FillPlacment(fillPlacement models.FillPlacement) (models.Placement, error) {
	seedWithBunkers, err := app.GetSeedWithBunkers(fillPlacement.Seed)
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

	return app.repo.PlcRepo.UpdatePlacement(placement)
}
