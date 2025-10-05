package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IBunkersRepository interface {
	AddBunkers(bunkers models.Bunkers) (bool, error)
	GetBunkers() ([]models.Bunkers, error)
	UpdateBunkers(bunkers models.Bunkers) (bool, error)
	DeleteBunkers(id int) (bool, error)
}

type bunkersRepository struct {
	db *sqlx.DB
}

func NewBunkersRepository(db *sqlx.DB) *bunkersRepository {
	return &bunkersRepository{
		db: db,
	}
}

func (bunk *bunkersRepository) AddBunkers(bunkers models.Bunkers) (bool, error) {
	const query = `
INSERT INTO green_seeds.bunkers (
	bunker,
	distance
)
VALUES (
	:bunker,
	:distance
)`

	result, err := bunk.db.NamedExec(query, bunkers)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (bunk *bunkersRepository) GetBunkers() ([]models.Bunkers, error) {
	const query = `
SELECT bunker, distance
FROM green_seeds.bunkers`

	var bunkers []models.Bunkers
	if err := bunk.db.Select(&bunkers, query); err != nil {
		return nil, err
	}

	return bunkers, nil
}

func (bunk *bunkersRepository) UpdateBunkers(bunkers models.Bunkers) (bool, error) {
	const query = `
UPDATE green_seeds.bunkers
SET distance = :distance
WHERE bunker = :bunker`

	result, err := bunk.db.NamedExec(query, bunkers)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
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
