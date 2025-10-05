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
// @Router /api/bunker/add [post]
// @Summary Добавление бункера в БД
// @Description При обращении, добавляет бункер в БД
//
// @Tags Bunkers
// @Produce      application/json
// @Consume      application/json
//
// @Param 	request	body	bunker	true	"Тело запроса"
//
// @Success 204 {object} nil "Запрос выполнен успешно"
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

	status, err := transport.service.AddBunker(bunker)
	if err != nil {
		utils.WriteJSON(w, status, err)
		return
	}

	utils.WriteNoContent(w)
}

// Set godoc
//
// @Router /api/bunker/get [get]
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
	bunkers, err := transport.service.GetBunkers()
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}

	if bunkers == nil {
		utils.WriteString(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, bunkers)
}

// Set godoc
//
// @Router /api/bunker/update [put]
// @Summary Обновление данных о бункере
// @Description При обращении, обновляет данные о бункере
//
// @Tags Bunkers
// @Produce      application/json
// @Consume      application/json
//
// @Success 204 {object} nil "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PutApiBunkerUpdate(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}

	var bunker models.Bunkers
	if err := jsoniter.Unmarshal(data, &bunker); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}

	status, err := transport.service.UpdateBunker(bunker)
	if err != nil {
		utils.WriteJSON(w, status, err)
		return
	}

	utils.WriteNoContent(w)
}

// Set godoc
//
// @Router /api/bunker/delete [delete]
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

	if err := transport.service.DeleteBunker(bunker); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteNoContent(w)
}
