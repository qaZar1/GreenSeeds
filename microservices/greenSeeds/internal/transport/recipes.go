package transport

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Impisigmatus/service_core/log"
	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/utils"
	"github.com/rs/zerolog"
)

// Set godoc
//
// @Router /api/recipes/add [post]
// @Summary Добавление информации, в каком бункере какие семена
// @Description При обращении, добавляет информацию о семенах в БД
//
// @Tags Recipes
// @Produce      application/json
// @Consume      application/json
//
// @Param 	request	body	recipes	true	"Тело запроса"
//
// @Success 200 {object} recipes "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiRecipesAdd(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %v", err))
		return
	}

	var recipe models.Recipes
	if err := jsoniter.Unmarshal(body, &recipe); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %v", err))
		return
	}

	addedRecipe, err := transport.Recipes.AddRecipes(recipe)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add recipes: %v", err))
		return
	}

	if addedRecipe == (models.Recipes{}) {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add placement: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, addedRecipe)
}

// Set godoc
//
// @Router /api/recipes/get [get]
// @Summary Получение списка бункеров и семян
// @Description При обращении, возвращает список бункеров и семян
//
// @Tags Recipes
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} []recipes "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiRecipesGet(w http.ResponseWriter, r *http.Request) {
	recipes, err := transport.Recipes.GetRecipes()
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get recipes: %v", err))
		return
	}

	if recipes == nil {
		utils.WriteString(w, http.StatusNotFound, fmt.Sprintf("Recipes not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, recipes)
}

// Set godoc
//
// @Router /api/recipes/get/{recipe} [get]
// @Summary Получение бункера и семян по ID
// @Description При обращении, возвращает бункер и семена по ID
//
// @Tags Recipes
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} recipes "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiRecipesGetRecipe(w http.ResponseWriter, r *http.Request) {
	recipeName := chi.URLParam(r, "recipe")

	recipeNum, err := strconv.Atoi(recipeName)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get recipe: %v", err))
		return
	}

	recipe, err := transport.Recipes.GetRecipesByRecipe(recipeNum)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get recipe: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, recipe)
}

// Set godoc
//
// @Router /api/recipes/update [put]
// @Summary Обновление данных о семенах
// @Description При обращении, обновляет данные о семенах
//
// @Tags Recipes
// @Produce      application/json
// @Consume      application/json
//
// @Param request body recipes true "Тело запроса"
//
// @Success 200 {object} recipes "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PutApiRecipesUpdate(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid read body: %v", err))
		return
	}

	var recipe models.Recipes
	if err := jsoniter.Unmarshal(data, &recipe); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid unmarshal: %v", err))
		return
	}

	updatedRecipe, err := transport.Recipes.UpdateRecipes(recipe)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update recipe: %v", err))
		return
	}

	if updatedRecipe == (models.Recipes{}) {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update recipe: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedRecipe)
}

// Set godoc
//
// @Router /api/recipes/delete/{recipe} [delete]
// @Summary Удаление рецепта
// @Description При обращении, удаляет рецепт
//
// @Tags Recipes
// @Produce      application/json
// @Consume      application/json
//
// @Param recipe path string true "Название рецепта"
//
// @Success 204 {object} nil "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) DeleteApiRecipesDelete(w http.ResponseWriter, r *http.Request) {
	log, ok := r.Context().Value(log.CtxKey).(zerolog.Logger)
	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid logger")
		return
	}

	recipeName := chi.URLParam(r, "recipe")

	recipeNum, err := strconv.Atoi(recipeName)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get recipe: %v", err))
		return
	}

	ok, err = transport.Recipes.DeleteRecipes(recipeNum)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid delete recipe")
		return
	}

	log.Warn().Ctx(r.Context()).Msg(fmt.Sprintf("Recipe removed: %s", recipeName))

	utils.WriteNoContent(w)
}
