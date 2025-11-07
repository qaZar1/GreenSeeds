package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type ISeedsRepository interface {
	AddSeeds(seeds models.Seeds) (models.Seeds, error)
	GetSeeds() ([]models.Seeds, error)
	GetSeedsBySeed(seedName string) (models.Seeds, error)
	UpdateSeeds(seeds models.Seeds) (models.Seeds, error)
	DeleteSeeds(seedName string) (bool, error)
}

type seedsRepository struct {
	db *sqlx.DB
}

func NewSeedsRepository(db *sqlx.DB) *seedsRepository {
	return &seedsRepository{
		db: db,
	}
}

func (se *seedsRepository) AddSeeds(seeds models.Seeds) (models.Seeds, error) {
	const query = `
INSERT INTO green_seeds.seeds (
	seed,
	seed_ru,
	min_density,
	max_density,
	tank_capacity,
	latency
)
VALUES (
	:seed,
	:seed_ru,
	:min_density,
	:max_density,
	:tank_capacity,
	:latency
)
RETURNING seed, seed_ru, min_density, max_density, tank_capacity, latency`

	rows, err := se.db.NamedQuery(query, seeds)
	if err != nil {
		return models.Seeds{}, err
	}

	defer rows.Close()

	var inserted models.Seeds
	if rows.Next() {
		err = rows.StructScan(&inserted)
		if err != nil {
			return models.Seeds{}, err
		}
	}

	return inserted, nil
}

func (se *seedsRepository) GetSeeds() ([]models.Seeds, error) {
	const query = `
SELECT seed, seed_ru, min_density, max_density, tank_capacity, latency
FROM green_seeds.seeds`

	var seeds []models.Seeds
	if err := se.db.Select(&seeds, query); err != nil {
		return nil, err
	}

	return seeds, nil
}

func (se *seedsRepository) GetSeedsBySeed(seedName string) (models.Seeds, error) {
	const query = `
SELECT seed, seed_ru, min_density, max_density, tank_capacity, latency
FROM green_seeds.seeds
WHERE seed = $1`

	var seed models.Seeds
	if err := se.db.Get(&seed, query, seedName); err != nil {
		return models.Seeds{}, err
	}

	return seed, nil
}

func (se *seedsRepository) UpdateSeeds(seeds models.Seeds) (models.Seeds, error) {
	const query = `
UPDATE green_seeds.seeds
SET	seed_ru = :seed_ru,
	min_density = :min_density,
	max_density = :max_density,
	tank_capacity = :tank_capacity,
	latency = :latency
WHERE seed = :seed
RETURNING seed, seed_ru, min_density, max_density, tank_capacity, latency`

	rows, err := se.db.NamedQuery(query, seeds)
	if err != nil {
		return models.Seeds{}, err
	}

	defer rows.Close()

	var updated models.Seeds
	if rows.Next() {
		err = rows.StructScan(&updated)
		if err != nil {
			return models.Seeds{}, err
		}
	}

	return updated, nil
}

func (se *seedsRepository) DeleteSeeds(seedName string) (bool, error) {
	const query = `
DELETE FROM green_seeds.seeds
WHERE seed = $1`

	result, err := se.db.Exec(query, seedName)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}
