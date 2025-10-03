package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IReceiptsRepository interface {
	AddReceipts(receipts *models.Receipts) (bool, error)
	GetReceipts() ([]models.Receipts, error)
	UpdateReceipts(receipts *models.Receipts) (bool, error)
	DeleteReceipts(receipt int) (bool, error)
}

type receiptsRepository struct {
	db *sqlx.DB
}

func NewReceiptsRepository(db *sqlx.DB) *receiptsRepository {
	return &receiptsRepository{
		db: db,
	}
}

func (rec *receiptsRepository) AddReceipts(receipts *models.Receipts) (bool, error) {
	const query = `
INSERT INTO green_seeds.receipts (
	receipt,
	seed,
	gcode,
	updated,
	description
)
VALUES (
	:receipt,
	:seed,
	:gcode,
	:updated,
	:description
)`

	result, err := rec.db.NamedExec(query, receipts)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (rec *receiptsRepository) GetReceipts() ([]models.Receipts, error) {
	const query = `
SELECT receipt, seed, gcode, updated, description
FROM green_seeds.receipts`

	var receipts []models.Receipts
	if err := rec.db.Select(&receipts, query); err != nil {
		return nil, err
	}

	return receipts, nil
}

func (rec *receiptsRepository) UpdateReceipts(receipts *models.Receipts) (bool, error) {
	const query = `
UPDATE green_seeds.receipts
SET seed = :seed,
	gcode = :gcode,
	updated = :updated,
	description = :description
WHERE receipt = :receipt`

	result, err := rec.db.NamedExec(query, receipts)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (rec *receiptsRepository) DeleteReceipts(receipt int) (bool, error) {
	const query = `
DELETE FROM green_seeds.receipts
WHERE receipt = $1`

	result, err := rec.db.Exec(query, receipt)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}
