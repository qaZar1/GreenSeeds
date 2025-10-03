package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type ISeedsRepository interface {
	AddSeeds(seeds *models.Seeds) (bool, error)
	GetSeeds() ([]models.Seeds, error)
	UpdateSeeds(seeds *models.Seeds) (bool, error)
	DeleteSeeds(seed string) (bool, error)
}

type seedsRepository struct {
	db *sqlx.DB
}

func NewSeedsRepository(db *sqlx.DB) *seedsRepository {
	return &seedsRepository{
		db: db,
	}
}

func (se *seedsRepository) AddSeeds(seeds *models.Seeds) (bool, error) {
	const query = `
INSERT INTO green_seeds.seeds (
	seed,
	min_density,
	max_density,
	tank_capacity,
	latency
)
VALUES (
	:seed,
	:min_density,
	:max_density,
	:tank_capacity,
	:latency
)`

	result, err := se.db.NamedExec(query, seeds)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (se *seedsRepository) GetSeeds() ([]models.Seeds, error) {
	const query = `
SELECT seed, min_density, max_density, tank_capacity, latency
FROM green_seeds.seeds`

	var seeds []models.Seeds
	if err := se.db.Select(&seeds, query); err != nil {
		return nil, err
	}

	return seeds, nil
}

func (se *seedsRepository) UpdateSeeds(seeds *models.Seeds) (bool, error) {
	const query = `
UPDATE green_seeds.seeds
SET	min_density = :min_density,
	max_density = :max_density,
	tank_capacity = :tank_capacity,
	latency = :latency
WHERE seed = :seed`

	result, err := se.db.NamedExec(query, seeds)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (se *seedsRepository) DeleteSeeds(seed string) (bool, error) {
	const query = `
DELETE FROM green_seeds.seeds
WHERE seed = $1`

	result, err := se.db.Exec(query, seed)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}
