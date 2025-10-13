package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IReceiptsRepository interface {
	AddReceipts(receipts models.Receipts) (models.Receipts, error)
	GetReceipts() ([]models.Receipts, error)
	UpdateReceipts(receipts models.Receipts) (models.Receipts, error)
	DeleteReceipts(receipt string) (bool, error)
	GetReceiptsByReceipt(receipt string) (models.Receipts, error)
}

type receiptsRepository struct {
	db *sqlx.DB
}

func NewReceiptsRepository(db *sqlx.DB) *receiptsRepository {
	return &receiptsRepository{
		db: db,
	}
}

func (rec *receiptsRepository) AddReceipts(receipts models.Receipts) (models.Receipts, error) {
	const query = `
INSERT INTO green_seeds.receipts (
	receipt,
	seed,
	gcode,
	description
)
VALUES (
	:receipt,
	:seed,
	:gcode,
	:description
)
RETURNING receipt, seed, gcode, description, updated`

	rows, err := rec.db.NamedQuery(query, receipts)
	if err != nil {
		return models.Receipts{}, err
	}

	defer rows.Close()

	var inserted models.Receipts
	if rows.Next() {
		err = rows.StructScan(&inserted)
		if err != nil {
			return models.Receipts{}, err
		}
	}

	return inserted, nil
}

func (rec *receiptsRepository) GetReceipts() ([]models.Receipts, error) {
	const query = `
SELECT receipt, seed, gcode, updated, description
FROM green_seeds.receipts
ORDER BY receipt ASC`

	var receipts []models.Receipts
	if err := rec.db.Select(&receipts, query); err != nil {
		return nil, err
	}

	return receipts, nil
}

func (rec *receiptsRepository) UpdateReceipts(receipts models.Receipts) (models.Receipts, error) {
	const query = `
UPDATE green_seeds.receipts
SET
	seed = :seed,
    gcode = :gcode,
    description = :description,
    updated = CURRENT_TIMESTAMP
WHERE receipt = :receipt
RETURNING receipt, seed, gcode, description, updated`

	rows, err := rec.db.NamedQuery(query, receipts)
	if err != nil {
		return models.Receipts{}, err
	}

	defer rows.Close()

	var updated models.Receipts
	if rows.Next() {
		err = rows.StructScan(&updated)
		if err != nil {
			return models.Receipts{}, err
		}
	}

	return updated, nil
}

func (rec *receiptsRepository) DeleteReceipts(receipt string) (bool, error) {
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

func (rec *receiptsRepository) GetReceiptsByReceipt(receiptName string) (models.Receipts, error) {
	const query = `
SELECT receipt, seed, gcode, updated, description
FROM green_seeds.receipts
WHERE receipt = $1`

	var receipt models.Receipts
	if err := rec.db.Get(&receipt, query, receiptName); err != nil {
		return models.Receipts{}, err
	}

	return receipt, nil
}
