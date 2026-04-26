package application

import (
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (app *App) AddSeed(seed models.Seeds) (models.Seeds, error) {
	if err := app.validate.Struct(seed); err != nil {
		return models.Seeds{}, err
	}

	return app.repo.SeedRepo.AddSeeds(seed)
}

func (app *App) GetSeeds() ([]models.Seeds, error) {
	return app.repo.SeedRepo.GetSeeds()
}

func (app *App) GetSeedBySeed(seed string) (models.Seeds, error) {
	return app.repo.SeedRepo.GetSeedsBySeed(seed)
}

func (app *App) GetSeedWithBunkers(seed string) ([]models.SeedsWithBunker, error) {
	return app.repo.SeedRepo.GetSeedsWithBunkers(seed)
}

func (app *App) UpdateSeed(seed models.Seeds) (models.Seeds, error) {
	if err := app.validate.Struct(seed); err != nil {
		return models.Seeds{}, err
	}

	return app.repo.SeedRepo.UpdateSeeds(seed)
}

func (app *App) DeleteSeed(seed string) (bool, error) {
	return app.repo.SeedRepo.DeleteSeeds(seed)
}
