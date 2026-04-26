package application

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (app *App) AddAssignment(assignment models.Assignments) (models.Assignments, error) {
	if err := app.validate.Struct(assignment); err != nil {
		return models.Assignments{}, err
	}

	return app.repo.AsnRepo.AddAssignments(assignment)
}

func (app *App) GetAssignments() ([]models.Assignments, error) {
	return app.repo.AsnRepo.GetAssignments()
}

func (app *App) GetAssignmentsByAssignment(idStr string) (models.Assignments, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.Assignments{}, err
	}

	return app.repo.AsnRepo.GetAssignmentsByNumber(id)
}

func (app *App) UpdateAssignment(assignment models.Assignments) (models.Assignments, error) {
	if err := app.validate.Struct(assignment); err != nil {
		return models.Assignments{}, err
	}

	if assignment.Id == nil {
		return models.Assignments{}, fmt.Errorf("Invalid ID")
	}

	oldAssignment, err := app.repo.AsnRepo.GetAssignmentsByNumber(int(*assignment.Id))
	if err != nil {
		return models.Assignments{}, err
	}

	reports, err := app.repo.RepRepo.GetReportsByAssignment(
		int(assignment.Shift),
		assignment.Number,
		int(assignment.Receipt))
	if err != nil {
		return models.Assignments{}, err
	}

	updated, err := app.repo.AsnRepo.SyncReports(assignment, oldAssignment, reports)
	if updated == (models.Assignments{}) || err != nil {
		return models.Assignments{}, errors.New("transaction failed")
	}

	return updated, nil
}

func (app *App) DeleteAssignments(idStr string) (bool, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return false, err
	}

	return app.repo.AsnRepo.DeleteAssignments(id)
}

func (app *App) CheckActiveTasks(userId string) ([]models.ActiveTask, error) {
	return app.repo.AsnRepo.CheckActiveTasks(userId)
}

func (app *App) GetTaskById(idStr string) (models.Task, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return models.Task{}, err
	}

	return app.repo.AsnRepo.GetTaskById(id)
}
