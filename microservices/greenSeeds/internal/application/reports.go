package application

import (
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (app *App) AddReport(report models.Reports) (models.Reports, error) {
	if err := app.validate.Struct(report); err != nil {
		return models.Reports{}, err
	}

	return app.repo.RepRepo.AddReports(report)
}

func (app *App) GetReports() ([]models.Reports, error) {
	return app.repo.RepRepo.GetReports()
}

func (app *App) GetReportsByReport(idStr string) (models.Reports, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.Reports{}, err
	}

	return app.repo.RepRepo.GetReportsById(id)
}

func (app *App) UpdateReport(report models.Reports) (bool, error) {
	if err := app.validate.Struct(report); err != nil {
		return false, err
	}

	return app.repo.RepRepo.UpdateReports(report)
}
