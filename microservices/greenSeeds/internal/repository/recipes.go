package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

type IRecipesRepository interface {
	AddRecipes(recipes models.Recipes) (models.Recipes, error)
	GetRecipes() ([]models.Recipes, error)
	UpdateRecipes(recipes models.Recipes) (models.Recipes, error)
	DeleteRecipes(recipe int) (bool, error)
	GetRecipesByRecipe(recipe int) (models.Recipes, error)
}

type recipesRepository struct {
	db *sqlx.DB
}

func NewRecipesRepository(db *sqlx.DB) *recipesRepository {
	return &recipesRepository{
		db: db,
	}
}

func (rec *recipesRepository) AddRecipes(recipes models.Recipes) (models.Recipes, error) {
	const query = `
INSERT INTO green_seeds.recipes (
	seed,
	gcode,
	description
)
VALUES (
	:seed,
	:gcode,
	:description
)
RETURNING recipe, seed, gcode, description, updated`

	rows, err := rec.db.NamedQuery(query, recipes)
	if err != nil {
		return models.Recipes{}, err
	}

	defer rows.Close()

	var inserted models.Recipes
	if rows.Next() {
		err = rows.StructScan(&inserted)
		if err != nil {
			return models.Recipes{}, err
		}
	}

	return inserted, nil
}

func (rec *recipesRepository) GetRecipes() ([]models.Recipes, error) {
	const query = `
SELECT 
    recipes.recipe, 
    recipes.seed, 
    seeds.seed_ru,
    recipes.gcode, 
    recipes.updated, 
    recipes.description
FROM green_seeds.recipes
LEFT JOIN green_seeds.seeds 
    ON seeds.seed = recipes.seed
WHERE recipes.deleted_at IS NULL
ORDER BY recipes.seed ASC;
`

	var recipes []models.Recipes
	if err := rec.db.Select(&recipes, query); err != nil {
		return nil, err
	}

	return recipes, nil
}

func (rec *recipesRepository) UpdateRecipes(recipes models.Recipes) (models.Recipes, error) {
	const query = `
UPDATE green_seeds.recipes
SET
	seed = :seed,
    gcode = :gcode,
    description = :description,
    updated = CURRENT_TIMESTAMP
WHERE recipe = :recipe AND deleted_at IS NULL
RETURNING recipe, seed, gcode, description, updated`

	rows, err := rec.db.NamedQuery(query, recipes)
	if err != nil {
		return models.Recipes{}, err
	}

	defer rows.Close()

	var updated models.Recipes
	if rows.Next() {
		err = rows.StructScan(&updated)
		if err != nil {
			return models.Recipes{}, err
		}
	}

	return updated, nil
}

func (rec *recipesRepository) DeleteRecipes(recipe int) (bool, error) {
	const query = `
UPDATE green_seeds.recipes
SET deleted_at = $1
WHERE recipe = $2 AND deleted_at IS NULL;`

	result, err := rec.db.Exec(query, time.Now(), recipe)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func (rec *recipesRepository) GetRecipesByRecipe(recipeNum int) (models.Recipes, error) {
	const query = `
SELECT 
    recipes.recipe, 
    recipes.seed, 
    seeds.seed_ru,
    recipes.gcode, 
    recipes.updated, 
    recipes.description
FROM green_seeds.recipes
LEFT JOIN green_seeds.seeds 
    ON seeds.seed = recipes.seed
WHERE recipe = $1 AND recipes.deleted_at IS NULL;`

	var recipe models.Recipes
	if err := rec.db.Get(&recipe, query, recipeNum); err != nil {
		return models.Recipes{}, err
	}

	return recipe, nil
}
