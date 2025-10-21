package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IBunkersRepository interface {
	AddBunkers(bunkers models.Bunkers) (models.Bunkers, error)
	GetBunkers() ([]models.Bunkers, error)
	GetBunkersById(id int) (models.Bunkers, error)
	UpdateBunkers(bunkers models.Bunkers) (models.Bunkers, error)
	DeleteBunkers(id int) (bool, error)
	GetBunkersForPlacement() ([]models.Bunkers, error)
}

type bunkersRepository struct {
	db *sqlx.DB
}

func NewBunkersRepository(db *sqlx.DB) *bunkersRepository {
	return &bunkersRepository{
		db: db,
	}
}

func (bunk *bunkersRepository) AddBunkers(bunkers models.Bunkers) (models.Bunkers, error) {
	const query = `
INSERT INTO green_seeds.bunkers (
	bunker,
	distance
)
VALUES (
	:bunker,
	:distance
)
RETURNING bunker, distance`

	rows, err := bunk.db.NamedQuery(query, bunkers)
	if err != nil {
		return models.Bunkers{}, err
	}

	defer rows.Close()

	var inserted models.Bunkers
	if rows.Next() {
		err = rows.StructScan(&inserted)
		if err != nil {
			return models.Bunkers{}, err
		}
	}

	return inserted, nil
}

func (bunk *bunkersRepository) GetBunkers() ([]models.Bunkers, error) {
	const query = `
SELECT bunker, distance
FROM green_seeds.bunkers
ORDER BY bunker ASC`

	var bunkers []models.Bunkers
	if err := bunk.db.Select(&bunkers, query); err != nil {
		return nil, err
	}

	return bunkers, nil
}

func (bunk *bunkersRepository) GetBunkersForPlacement() ([]models.Bunkers, error) {
	const query = `
SELECT bunker, distance
FROM green_seeds.bunkers
WHERE bunker NOT IN (SELECT bunker FROM green_seeds.placement)
ORDER BY bunker ASC`

	var bunkers []models.Bunkers
	if err := bunk.db.Select(&bunkers, query); err != nil {
		return nil, err
	}

	return bunkers, nil
}

func (bunk *bunkersRepository) GetBunkersById(id int) (models.Bunkers, error) {
	const query = `
SELECT bunker, distance
FROM green_seeds.bunkers
WHERE bunker = $1`

	var bunker models.Bunkers
	if err := bunk.db.Get(&bunker, query, id); err != nil {
		return models.Bunkers{}, err
	}

	return bunker, nil
}

func (bunk *bunkersRepository) UpdateBunkers(bunkers models.Bunkers) (models.Bunkers, error) {
	const query = `
UPDATE green_seeds.bunkers
SET distance = :distance
WHERE bunker = :bunker
RETURNING bunker, distance`

	rows, err := bunk.db.NamedQuery(query, bunkers)
	if err != nil {
		return models.Bunkers{}, err
	}

	defer rows.Close()

	var updated models.Bunkers
	if rows.Next() {
		err = rows.StructScan(&updated)
		if err != nil {
			return models.Bunkers{}, err
		}
	}

	return updated, nil
}

func (bunk *bunkersRepository) DeleteBunkers(bunker int) (bool, error) {
	const query = `
DELETE FROM green_seeds.bunkers
WHERE bunker = $1`

	result, err := bunk.db.Exec(query, bunker)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}
