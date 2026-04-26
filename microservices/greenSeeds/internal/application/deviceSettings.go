package application

import (
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (app *App) AddSetting(setting models.DeviceSettings) (models.DeviceSettings, error) {
	if err := app.validate.Struct(setting); err != nil {
		return models.DeviceSettings{}, err
	}

	return app.repo.DevSet.AddSetting(setting)
}

func (app *App) GetSettings() ([]models.DeviceSettings, error) {
	return app.repo.DevSet.GetSettings()
}

func (app *App) GetSettingsByKey(key string) (models.DeviceSettings, error) {
	return app.repo.DevSet.GetSettingsByKey(key)
}

func (app *App) UpdateSetting(setting models.DeviceSettings) (models.DeviceSettings, error) {
	if err := app.validate.Struct(setting); err != nil {
		return models.DeviceSettings{}, err
	}

	return app.repo.DevSet.UpdateSettings(setting)
}

func (app *App) DeleteSetting(key string) (bool, error) {
	return app.repo.DevSet.DeleteSettings(key)
}
