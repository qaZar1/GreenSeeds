package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IReceiptsRepository interface {
	AddReceipts(receipts models.Receipts) (models.Receipts, error)
	GetReceipts() ([]models.Receipts, error)
	UpdateReceipts(receipts models.Receipts) (models.Receipts, error)
	DeleteReceipts(receipt int) (bool, error)
	GetReceiptsByReceipt(receipt int) (models.Receipts, error)
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
	seed,
	gcode,
	description
)
VALUES (
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
SELECT 
    receipts.receipt, 
    receipts.seed, 
    seeds.seed_ru,
    receipts.gcode, 
    receipts.updated, 
    receipts.description
FROM green_seeds.receipts
LEFT JOIN green_seeds.seeds 
    ON seeds.seed = receipts.seed
WHERE receipts.deleted_at IS NULL
ORDER BY receipts.seed ASC;
`

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
WHERE receipt = :receipt AND deleted_at IS NULL
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

func (rec *receiptsRepository) DeleteReceipts(receipt int) (bool, error) {
	const query = `
UPDATE green_seeds.receipts
SET deleted_at = $1
WHERE receipt = $2 AND deleted_at IS NULL;`

	result, err := rec.db.Exec(query, time.Now(), receipt)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (rec *receiptsRepository) GetReceiptsByReceipt(receiptNum int) (models.Receipts, error) {
	const query = `
SELECT 
    receipts.receipt, 
    receipts.seed, 
    seeds.seed_ru,
    receipts.gcode, 
    receipts.updated, 
    receipts.description
FROM green_seeds.receipts
LEFT JOIN green_seeds.seeds 
    ON seeds.seed = receipts.seed
WHERE receipt = $1 AND receipts.deleted_at IS NULL;`

	var receipt models.Receipts
	if err := rec.db.Get(&receipt, query, receiptNum); err != nil {
		return models.Receipts{}, err
	}

	return receipt, nil
}
