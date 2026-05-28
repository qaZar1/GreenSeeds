package application

import (
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

//go:generate mockgen -source=recipes.go -destination=./../mocks/mock_recipes.go -package=mocks
type IRecipesApp interface {
	AddRecipes(models.Recipes) (models.Recipes, error)
	GetRecipes() ([]models.Recipes, error)
	GetRecipesByRecipe(int) (models.Recipes, error)
	UpdateRecipes(models.Recipes) (models.Recipes, error)
	DeleteRecipes(int) (bool, error)
}

func (app *App) AddRecipes(recipes models.Recipes) (models.Recipes, error) {
	if err := app.validate.Struct(recipes); err != nil {
		return models.Recipes{}, err
	}

	return app.repo.RptRepo.AddRecipes(recipes)
}

func (app *App) GetRecipes() ([]models.Recipes, error) {
	return app.repo.RptRepo.GetRecipes()
}

func (app *App) GetRecipesByRecipe(recipe int) (models.Recipes, error) {
	return app.repo.RptRepo.GetRecipesByRecipe(recipe)
}

func (app *App) UpdateRecipes(recipes models.Recipes) (models.Recipes, error) {
	if err := app.validate.Struct(recipes); err != nil {
		return models.Recipes{}, err
	}

	return app.repo.RptRepo.UpdateRecipes(recipes)
}

func (app *App) DeleteRecipes(recipe int) (bool, error) {
	return app.repo.RptRepo.DeleteRecipes(recipe)
}
