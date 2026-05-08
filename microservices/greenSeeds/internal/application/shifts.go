package application

import (
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

//go:generate mockgen -source=shifts.go -destination=./../mocks/mock_shifts.go -package=mocks
type IShiftsApp interface {
	AddShift(models.Shifts) (models.Shifts, error)
	GetShifts() ([]models.Shifts, error)
	GetShiftsByShift(string) (models.Shifts, error)
	UpdateShifts(models.Shifts) (models.Shifts, error)
	DeleteShifts(string) (bool, error)
	GetShiftsWithoutUser() ([]models.Shifts, error)
}

func (app *App) AddShift(shifts models.Shifts) (models.Shifts, error) {
	if err := app.validate.Struct(shifts); err != nil {
		return models.Shifts{}, err
	}

	return app.repo.ShfRepo.AddShifts(shifts)
}

func (app *App) GetShifts() ([]models.Shifts, error) {
	return app.repo.ShfRepo.GetShifts()
}

func (app *App) GetShiftsByShift(shift string) (models.Shifts, error) {
	shiftInt, err := strconv.Atoi(shift)
	if err != nil {
		return models.Shifts{}, err
	}

	return app.repo.ShfRepo.GetShiftsByShift(shiftInt)
}

func (app *App) UpdateShifts(shifts models.Shifts) (models.Shifts, error) {
	if err := app.validate.Struct(shifts); err != nil {
		return models.Shifts{}, err
	}

	return app.repo.ShfRepo.UpdateShifts(shifts)
}

func (app *App) DeleteShifts(shift string) (bool, error) {
	shiftInt, err := strconv.Atoi(shift)
	if err != nil {
		return false, err
	}

	return app.repo.ShfRepo.DeleteShifts(shiftInt)
}

func (app *App) GetShiftsWithoutUser() ([]models.Shifts, error) {
	return app.repo.ShfRepo.GetShiftsWithoutUser()
}
