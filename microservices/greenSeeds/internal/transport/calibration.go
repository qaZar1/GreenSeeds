package transport

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Impisigmatus/service_core/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/utils"
	"github.com/rs/zerolog"
)

// Set godoc
//
// @Router /api/calibration/handshake [post]
// @Summary Проверка подключения к устройству
// @Description При обращении, проверяет подключение к устройству
//
// @Tags Calibration
// @Produce      application/json
// @Consume      application/json
//
// @Param 	request	body	WSMessage	true	"Тело запроса"
//
// @Success 200 {object} WSMessage "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiCalibrationHandshake(w http.ResponseWriter, r *http.Request) {
	response, err := transport.service.CalibrationHandshake()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("Invalid calibration handshake: %s", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// Set godoc
//
// @Router /api/calibration/photo [post]
// @Summary Запрос фотографии
// @Description При обращении, делает фотографию
//
// @Tags Calibration
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} models.Calibration "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiCalibrationPhoto(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %w", err))
		return
	}

	var calibration models.Calculation
	if err := jsoniter.Unmarshal(body, &calibration); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %w", err))
		return
	}

	photo, err := transport.service.GetPhoto(calibration)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get photo: %s", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, photo)
}

// Set godoc
//
// @Router /api/calibration/clear [post]
// @Summary Получение данных о задании на смену
// @Description При обращении, возвращает данные о задании на смену
//
// @Tags Assignments
// @Produce      application/json
// @Consume      application/json
//
// @Params id path int true "ID задания на смену"
//
// @Success 200 {object} assignment "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiCalibrationClear(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %w", err))
		return
	}

	var calibration models.Calculation
	if err := jsoniter.Unmarshal(body, &calibration); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %w", err))
		return
	}

	if err := transport.service.Clear(calibration.SessionId); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid clear: %s", err))
		return
	}

	utils.WriteNoContent(w)
}

// Set godoc
//
// @Router /api/calibration/found [post]
// @Summary Обновление данных о задании на смену
// @Description При обращении, обновляет данные о задании на смену
//
// @Tags Assignments
// @Produce      application/json
// @Consume      application/json
//
// @Param request body calibration true "Тело запроса"
//
// @Success 200 {object} calibration "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiCalibrationFound(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid read body: %s", err))
		return
	}

	var calibration models.Calibration
	if err := jsoniter.Unmarshal(data, &calibration); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid unmarshal: %s", err))
		return
	}

	found, err := transport.service.CalculateResult(calibration)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update assignment: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, found)
}

// Set godoc
//
// @Router /api/calibration/save [post]
// @Summary Обновление данных о задании на смену
// @Description При обращении, обновляет данные о задании на смену
//
// @Tags Assignments
// @Produce      application/json
// @Consume      application/json
//
// @Param request body models.Distance true "Тело запроса"
//
// @Success 200 {object} calibration "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiCalibrationSave(w http.ResponseWriter, r *http.Request) {
	log, ok := r.Context().Value(log.CtxKey).(zerolog.Logger)
	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid logger")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid read body: %s", err))
		return
	}

	var distance models.Calibration
	if err := jsoniter.Unmarshal(data, &distance); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid unmarshal: %s", err))
		return
	}

	if err := transport.service.Save(distance); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update distance by step: %s", err))
		return
	}

	log.Warn().Ctx(r.Context()).Msg(fmt.Sprintf("Calibration saved %s", distance.SessionId))

	utils.WriteNoContent(w)
}
