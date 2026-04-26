package application

import (
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (app *App) GetLogs(params models.LogsParams) ([]models.Log, error) {
	return app.repo.LogsRepo.GetLogs(params)
}
