package transport

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/utils"
)

// Set godoc
//
// @Router /api/seeds/add [post]
// @Summary Добавление информации о семенах в БД
// @Description При обращении, добавляет информацию о семенах в БД
//
// @Tags Seeds
// @Produce      application/json
// @Consume      application/json
//
// @Param 	request	body	seed	true	"Тело запроса"
//
// @Success 200 {object} seed "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiSeedAdd(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %w", err))
		return
	}

	var seed models.Seeds
	if err := jsoniter.Unmarshal(body, &seed); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %w", err))
		return
	}

	ok, err := transport.service.AddSeed(seed)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add seed: %w", err))
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add seed: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, seed)
}

// Set godoc
//
// @Router /api/seeds/get [get]
// @Summary Получение списка семян
// @Description При обращении, возвращает список семян
//
// @Tags Seeds
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} []seed "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiSeedGet(w http.ResponseWriter, r *http.Request) {
	seeds, err := transport.service.GetSeeds()
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get seeds: %w", err))
		return
	}

	if seeds == nil {
		utils.WriteString(w, http.StatusNotFound, fmt.Sprintf("Seeds not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, seeds)
}

// Set godoc
//
// @Router /api/seeds/get/{seed} [get]
// @Summary Получение семян по ID
// @Description При обращении, возвращает семена по ID
//
// @Tags Seeds
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} seed "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiSeedGetId(w http.ResponseWriter, r *http.Request) {
	seedId := chi.URLParam(r, "seed")

	seed, err := transport.service.GetSeedById(seedId)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get seed by id: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, seed)
}

// Set godoc
//
// @Router /api/seeds/update [put]
// @Summary Обновление данных о семенах
// @Description При обращении, обновляет данные о семенах
//
// @Tags Seeds
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} seed "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PutApiSeedUpdate(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid read body: %w", err))
		return
	}

	var seed models.Seeds
	if err := jsoniter.Unmarshal(data, &seed); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid unmarshal: %w", err))
		return
	}

	ok, err := transport.service.UpdateSeed(seed)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update seed: %w", err))
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update seed: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, seed)
}

// Set godoc
//
// @Router /api/seeds/delete/{seed} [delete]
// @Summary Удаление семян
// @Description При обращении, удаляет семена
//
// @Tags Seeds
// @Produce      application/json
// @Consume      application/json
//
// @Param seed path int true "ID семян"
//
// @Success 204 {object} nil "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) DeleteApiSeedDelete(w http.ResponseWriter, r *http.Request) {
	seedId := chi.URLParam(r, "seed")

	ok, err := transport.service.DeleteSeed(seedId)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid delete seed: %w", err))
		return
	}

	utils.WriteNoContent(w)
}
