package transport

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Impisigmatus/service_core/log"
	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/utils"
	"github.com/rs/zerolog"
)

// Set godoc
//
// @Router /api/shifts/add [post]
// @Summary Добавление информации о смене
// @Description При обращении, добавляет информацию о смене в БД
//
// @Tags Shifts
// @Produce      application/json
// @Consume      application/json
//
// @Param 	request	body	shift	true	"Тело запроса"
//
// @Success 200 {object} shift "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiShiftAdd(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %w", err))
		return
	}

	var shift models.Shifts
	if err := jsoniter.Unmarshal(body, &shift); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %w", err))
		return
	}

	addedShift, err := transport.app.AddShift(shift)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add shift: %w", err))
		return
	}

	if addedShift == (models.Shifts{}) {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add shift: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, addedShift)
}

// Set godoc
//
// @Router /api/shifts/get [get]
// @Summary Получение списка смен
// @Description При обращении, возвращает список смен
//
// @Tags Shifts
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} []shift "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiShiftsGet(w http.ResponseWriter, r *http.Request) {
	shifts, err := transport.app.GetShifts()
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get shifts: %w", err))
		return
	}

	if shifts == nil {
		utils.WriteString(w, http.StatusNotFound, fmt.Sprintf("Shifts not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, shifts)
}

// Set godoc
//
// @Router /api/shifts/get/{shift} [get]
// @Summary Получение данных о смене
// @Description При обращении, возвращает данные о смене
//
// @Tags Shifts
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} shift "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiShiftsGetShift(w http.ResponseWriter, r *http.Request) {
	shiftName := chi.URLParam(r, "shift")

	shift, err := transport.app.GetShiftsByShift(shiftName)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get shift by shift: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, shift)
}

// Set godoc
//
// @Router /api/shifts/update [put]
// @Summary Обновление данных о смене
// @Description При обращении, обновляет данные о смене
//
// @Tags Shifts
// @Produce      application/json
// @Consume      application/json
//
// @Param request body shift true "Тело запроса"
//
// @Success 200 {object} shift "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PutApiShiftsUpdate(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid read body: %w", err))
		return
	}

	var shift models.Shifts
	if err := jsoniter.Unmarshal(data, &shift); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid unmarshal: %w", err))
		return
	}

	updatedShift, err := transport.app.UpdateShifts(shift)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update shift: %w", err))
		return
	}

	if updatedShift == (models.Shifts{}) {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update shift: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedShift)
}

// Set godoc
//
// @Router /api/shifts/delete/{shift} [delete]
// @Summary Удаление смены
// @Description При обращении, удаляет смену
//
// @Tags Shifts
// @Produce      application/json
// @Consume      application/json
//
// @Param shift path string true "Название смены"
//
// @Success 204 {object} nil "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) DeleteApiShiftsDelete(w http.ResponseWriter, r *http.Request) {
	log, ok := r.Context().Value(log.CtxKey).(zerolog.Logger)
	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid logger")
		return
	}

	shiftName := chi.URLParam(r, "shift")

	ok, err := transport.app.DeleteShifts(shiftName)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid delete shift")
		return
	}

	log.Warn().Ctx(r.Context()).Msg(fmt.Sprintf("Shift removed: %s", shiftName))

	utils.WriteNoContent(w)
}

// Set godoc
//
// @Router /api/shifts/getWithoutUser [get]
// @Summary Получение списка смен
// @Description При обращении, возвращает список смен
//
// @Tags Shifts
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} []shift "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiShiftsGetWithoutUser(w http.ResponseWriter, r *http.Request) {
	shifts, err := transport.app.GetShiftsWithoutUser()
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get shifts: %w", err))
		return
	}

	if shifts == nil {
		utils.WriteString(w, http.StatusNotFound, fmt.Sprintf("Shifts not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, shifts)
}
