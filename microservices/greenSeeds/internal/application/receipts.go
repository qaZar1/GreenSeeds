package application

import (
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

//go:generate mockgen -source=receipts.go -destination=./../mocks/mock_receipts.go -package=mocks
type IReceiptsApp interface {
	AddReceipts(models.Receipts) (models.Receipts, error)
	GetReceipts() ([]models.Receipts, error)
	GetReceiptsByReceipt(int) (models.Receipts, error)
	UpdateReceipts(models.Receipts) (models.Receipts, error)
	DeleteReceipts(int) (bool, error)
}

func (app *App) AddReceipts(receipts models.Receipts) (models.Receipts, error) {
	if err := app.validate.Struct(receipts); err != nil {
		return models.Receipts{}, err
	}

	return app.repo.RptRepo.AddReceipts(receipts)
}

func (app *App) GetReceipts() ([]models.Receipts, error) {
	return app.repo.RptRepo.GetReceipts()
}

func (app *App) GetReceiptsByReceipt(receipt int) (models.Receipts, error) {
	return app.repo.RptRepo.GetReceiptsByReceipt(receipt)
}

func (app *App) UpdateReceipts(receipts models.Receipts) (models.Receipts, error) {
	if err := app.validate.Struct(receipts); err != nil {
		return models.Receipts{}, err
	}

	return app.repo.RptRepo.UpdateReceipts(receipts)
}

func (app *App) DeleteReceipts(receipt int) (bool, error) {
	return app.repo.RptRepo.DeleteReceipts(receipt)
}
