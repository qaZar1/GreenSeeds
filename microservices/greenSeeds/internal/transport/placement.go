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
// @Router /api/placement/add [post]
// @Summary Добавление информации, в каком бункере какие семена
// @Description При обращении, добавляет информацию о семенах в БД
//
// @Tags Placements
// @Produce      application/json
// @Consume      application/json
//
// @Param 	request	body	placement	true	"Тело запроса"
//
// @Success 200 {object} placement "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiPlacementAdd(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %w", err))
		return
	}

	var placement models.Placement
	if err := jsoniter.Unmarshal(body, &placement); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %w", err))
		return
	}

	ok, err := transport.service.AddPlacement(placement)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add placement: %w", err))
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add placement: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, placement)
}

// Set godoc
//
// @Router /api/placement/get [get]
// @Summary Получение списка бункеров и семян
// @Description При обращении, возвращает список бункеров и семян
//
// @Tags Placements
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} []placement "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiPlacementGet(w http.ResponseWriter, r *http.Request) {
	placements, err := transport.service.GetPlacements()
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get placements: %w", err))
		return
	}

	if placements == nil {
		utils.WriteString(w, http.StatusNotFound, fmt.Sprintf("Placements not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, placements)
}

// Set godoc
//
// @Router /api/placement/get/{bunker} [get]
// @Summary Получение бункера и семян по ID
// @Description При обращении, возвращает бункер и семена по ID
//
// @Tags Placements
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} placement "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiPlacementGetBunker(w http.ResponseWriter, r *http.Request) {
	bunkerId := chi.URLParam(r, "bunker")

	placement, err := transport.service.GetPlacementByBunker(bunkerId)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get placement by bunker: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, placement)
}

// Set godoc
//
// @Router /api/placement/update [put]
// @Summary Обновление данных о семенах
// @Description При обращении, обновляет данные о семенах
//
// @Tags Placements
// @Produce      application/json
// @Consume      application/json
//
// @Param request body placement true "Тело запроса"
//
// @Success 200 {object} placement "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PutApiPlacementUpdate(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid read body: %w", err))
		return
	}

	var placement models.Placement
	if err := jsoniter.Unmarshal(data, &placement); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid unmarshal: %w", err))
		return
	}

	ok, err := transport.service.UpdatePlacement(placement)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update placement: %w", err))
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update placement: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, placement)
}

// Set godoc
//
// @Router /api/placement/delete/{bunker} [delete]
// @Summary Удаление бункера
// @Description При обращении, удаляет бункер
//
// @Tags Placements
// @Produce      application/json
// @Consume      application/json
//
// @Param bunker path int true "ID бункера"
//
// @Success 204 {object} nil "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) DeleteApiPlacementDelete(w http.ResponseWriter, r *http.Request) {
	bunkerId := chi.URLParam(r, "bunker")

	ok, err := transport.service.DeletePlacement(bunkerId)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid delete placement")
		return
	}

	utils.WriteNoContent(w)
}
