package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type ISeedsRepository interface {
	AddSeeds(seeds models.Seeds) (models.Seeds, error)
	GetSeeds() ([]models.Seeds, error)
	GetSeedsBySeed(seedName string) (models.Seeds, error)
	UpdateSeeds(seeds models.Seeds) (models.Seeds, error)
	DeleteSeeds(seedName string) (bool, error)
	GetSeedsWithBunkers(seed string) ([]models.SeedsWithBunker, error)
	GetBestBunker(seed string) (models.SeedsWithBunker, error)
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
	tank_capacity
)
VALUES (
	:seed,
	:seed_ru,
	:min_density,
	:max_density,
	:tank_capacity
)
RETURNING seed, seed_ru, min_density, max_density, tank_capacity`

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
SELECT seed, seed_ru, min_density, max_density, tank_capacity, deleted_at
FROM green_seeds.seeds;`

	var seeds []models.Seeds
	if err := se.db.Select(&seeds, query); err != nil {
		return nil, err
	}

	return seeds, nil
}

func (se *seedsRepository) GetSeedsBySeed(seedName string) (models.Seeds, error) {
	const query = `
SELECT seed, seed_ru, min_density, max_density, tank_capacity
FROM green_seeds.seeds
WHERE seed = $1 AND deleted_at IS NULL;`

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
	deleted_at = :deleted_at
WHERE seed = :seed
RETURNING seed, seed_ru, min_density, max_density, tank_capacity, deleted_at`

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
UPDATE green_seeds.seeds
SET deleted_at = $1
WHERE seed = $2 AND deleted_at IS NULL;`

	result, err := se.db.Exec(query, time.Now(), seedName)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (se *seedsRepository) GetSeedsWithBunkers(seed string) ([]models.SeedsWithBunker, error) {
	const query = `
SELECT 
	s.seed_ru,
	s.seed,
	s.min_density,
	s.max_density,
	s.tank_capacity,
	p.amount,
	b.bunker
FROM green_seeds.seeds s
LEFT JOIN green_seeds.placement p
ON p.seed = s.seed
LEFT JOIN green_seeds.bunkers b 
ON p.bunker = b.bunker
WHERE s.seed = $1 AND s.deleted_at IS NULL;`

	var seeds []models.SeedsWithBunker
	if err := se.db.Select(&seeds, query, seed); err != nil {
		return nil, err
	}

	return seeds, nil
}

func (se *seedsRepository) GetBestBunker(seed string) (models.SeedsWithBunker, error) {
	const query = `
SELECT 
	s.seed_ru,
	s.seed,
	s.min_density,
	s.max_density,
	s.tank_capacity,
	p.amount,
	b.bunker
FROM green_seeds.seeds s
LEFT JOIN green_seeds.placement p
	ON p.seed = s.seed
LEFT JOIN green_seeds.bunkers b 
	ON p.bunker = b.bunker
WHERE s.seed = $1 
	AND s.deleted_at IS NULL
	AND p.amount > 0
ORDER BY p.amount DESC
LIMIT 1;`

	var seeds models.SeedsWithBunker
	if err := se.db.Get(&seeds, query, seed); err != nil {
		return models.SeedsWithBunker{}, err
	}

	return seeds, nil
}

