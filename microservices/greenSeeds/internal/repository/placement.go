package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IPlacementRepository interface {
	AddPlacement(placement models.Placement) (models.Placement, error)
	GetPlacement() ([]models.Placement, error)
	UpdatePlacement(placement models.Placement) (models.Placement, error)
	DeletePlacement(bunker int) (bool, error)
	GetPlacementByBunker(bunker int) (models.Placement, error)
}

type placementRepository struct {
	db *sqlx.DB
}

func NewPlacementRepository(db *sqlx.DB) *placementRepository {
	return &placementRepository{
		db: db,
	}
}

func (pl *placementRepository) AddPlacement(placement models.Placement) (models.Placement, error) {
	const query = `
INSERT INTO green_seeds.placement (
	bunker,
	seed
)
VALUES (
	:bunker,
	:seed
)
RETURNING bunker, seed`

	rows, err := pl.db.NamedQuery(query, placement)
	if err != nil {
		return models.Placement{}, err
	}

	defer rows.Close()

	var inserted models.Placement
	if rows.Next() {
		err = rows.StructScan(&inserted)
		if err != nil {
			return models.Placement{}, err
		}
	}

	return inserted, nil
}

func (pl *placementRepository) GetPlacement() ([]models.Placement, error) {
	const query = `
SELECT bunker, seed
FROM green_seeds.placement
ORDER BY bunker ASC`

	var placement []models.Placement
	if err := pl.db.Select(&placement, query); err != nil {
		return nil, err
	}

	return placement, nil
}

func (pl *placementRepository) UpdatePlacement(placement models.Placement) (models.Placement, error) {
	const query = `
UPDATE green_seeds.placement
SET	seed = :seed
WHERE bunker = :bunker
RETURNING bunker, seed`

	rows, err := pl.db.NamedQuery(query, placement)
	if err != nil {
		return models.Placement{}, err
	}

	defer rows.Close()

	var updated models.Placement
	if rows.Next() {
		err = rows.StructScan(&updated)
		if err != nil {
			return models.Placement{}, err
		}
	}

	return updated, nil
}

func (pl *placementRepository) DeletePlacement(bunker int) (bool, error) {
	const query = `
DELETE FROM green_seeds.placement
WHERE bunker = $1`

	result, err := pl.db.Exec(query, bunker)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (pl *placementRepository) GetPlacementByBunker(bunker int) (models.Placement, error) {
	const query = `
SELECT bunker, seed
FROM green_seeds.placement
WHERE bunker = $1`

	var placement models.Placement
	if err := pl.db.Get(&placement, query, bunker); err != nil {
		return models.Placement{}, err
	}

	return placement, nil
}
