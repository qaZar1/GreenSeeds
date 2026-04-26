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
// @Router /api/bunkers/add [post]
// @Summary Добавление бункера в БД
// @Description При обращении, добавляет бункер в БД
//
// @Tags Bunkers
// @Produce      application/json
// @Consume      application/json
//
// @Param 	request	body	bunker	true	"Тело запроса"
//
// @Success 200 {object} bunker "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiBunkerAdd(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %w", err))
		return
	}

	var bunker models.Bunkers
	if err := jsoniter.Unmarshal(body, &bunker); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %w", err))
		return
	}

	inserted, err := transport.app.AddBunker(bunker)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add bunker: %w", err))
		return
	}

	if inserted == (models.Bunkers{}) {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add bunker: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, inserted)
}

// Set godoc
//
// @Router /api/bunkers/get [get]
// @Summary Получение списка бункеров
// @Description При обращении, возвращает список бункеров
//
// @Tags Bunkers
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} []bunker "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiBunkerGet(w http.ResponseWriter, r *http.Request) {
	bunkers, err := transport.app.GetBunkers()
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get bunkers: %w", err))
		return
	}

	if bunkers == nil {
		utils.WriteString(w, http.StatusNotFound, fmt.Sprintf("Bunkers not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, bunkers)
}

// Set godoc
//
// @Router /api/bunkers/getForPlacement [get]
// @Summary Получение списка бункеров для размещения
// @Description При обращении, возвращает список бункеров для размещения
//
// @Tags Bunkers
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} []bunker "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiBunkerGetForPlacement(w http.ResponseWriter, r *http.Request) {
	bunkers, err := transport.app.GetBunkersForPlacement()
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get bunkers for placement: %w", err))
		return
	}

	if bunkers == nil {
		utils.WriteString(w, http.StatusNotFound, fmt.Sprintf("Bunkers not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, bunkers)
}

// Set godoc
//
// @Router /api/bunkers/get/{bunker} [get]
// @Summary Получение бункера по ID
// @Description При обращении, возвращает бункер по ID
//
// @Tags Bunkers
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} bunker "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiBunkerGetId(w http.ResponseWriter, r *http.Request) {
	bunkerId := chi.URLParam(r, "bunker")

	bunker, err := transport.app.GetBunkersById(bunkerId)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get bunker by id: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, bunker)
}

// Set godoc
//
// @Router /api/bunkers/update [put]
// @Summary Обновление данных о бункере
// @Description При обращении, обновляет данные о бункере
//
// @Tags Bunkers
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} bunker "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PutApiBunkerUpdate(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid read body: %w", err))
		return
	}

	var bunker models.Bunkers
	if err := jsoniter.Unmarshal(data, &bunker); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid unmarshal: %w", err))
		return
	}

	updated, err := transport.app.UpdateBunker(bunker)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update bunker: %w", err))
		return
	}

	if updated == (models.Bunkers{}) {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update bunker: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, updated)
}

// Set godoc
//
// @Router /api/bunkers/delete/{bunker} [delete]
// @Summary Удаление бункера
// @Description При обращении, удаляет бункер
//
// @Tags Bunkers
// @Produce      application/json
// @Consume      application/json
//
// @Param bunker path int true "ID бункера"
//
// @Success 204 {object} nil "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) DeleteApiBunkerDelete(w http.ResponseWriter, r *http.Request) {
	bunker := chi.URLParam(r, "bunker")

	ok, err := transport.app.DeleteBunker(bunker)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid delete bunker: %w", err))
		return
	}

	utils.WriteNoContent(w)
}
