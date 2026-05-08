package application

import (
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

//go:generate mockgen -source=logs.go -destination=./../mocks/mock_logs.go -package=mocks
type ILogsApp interface {
	GetLogs(models.LogsParams) ([]models.Log, error)
}

func (app *App) GetLogs(params models.LogsParams) ([]models.Log, error) {
	return app.repo.LogsRepo.GetLogs(params)
}
