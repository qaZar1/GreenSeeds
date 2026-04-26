package application

import (
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (app *App) AddBunker(bunker models.Bunkers) (models.Bunkers, error) {
	if err := app.validate.Struct(bunker); err != nil {
		return models.Bunkers{}, err
	}

	return app.repo.BunkRepo.AddBunkers(bunker)
}

func (app *App) GetBunkers() ([]models.Bunkers, error) {
	return app.repo.BunkRepo.GetBunkers()
}

func (app *App) GetBunkersForPlacement() ([]models.Bunkers, error) {
	return app.repo.BunkRepo.GetBunkersForPlacement()
}

func (app *App) GetBunkersById(bunkerId string) (models.Bunkers, error) {
	bunkerIdInt, err := strconv.Atoi(bunkerId)
	if err != nil {
		return models.Bunkers{}, err
	}

	return app.repo.BunkRepo.GetBunkersById(bunkerIdInt)
}

func (app *App) UpdateBunker(bunker models.Bunkers) (models.Bunkers, error) {
	if err := app.validate.Struct(bunker); err != nil {
		return models.Bunkers{}, err
	}

	return app.repo.BunkRepo.UpdateBunkers(bunker)
}

func (app *App) DeleteBunker(bunkerId string) (bool, error) {
	bunkerIdInt, err := strconv.Atoi(bunkerId)
	if err != nil {
		return false, err
	}

	return app.repo.BunkRepo.DeleteBunkers(bunkerIdInt)
}
