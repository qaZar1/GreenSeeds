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
// @Router /api/reports/add [post]
// @Summary Добавление информации о задании на смену
// @Description При обращении, добавляет информацию о задании на смену в БД
//
// @Tags Reports
// @Produce      application/json
// @Consume      application/json
//
// @Param 	request	body	assignment	true	"Тело запроса"
//
// @Success 200 {object} assignment "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiReportsAdd(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %w", err))
		return
	}

	var report models.Reports
	if err := jsoniter.Unmarshal(body, &report); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %w", err))
		return
	}

	addedReport, err := transport.service.AddReport(report)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add report: %w", err))
		return
	}

	if addedReport == (models.Reports{}) {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add report: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, addedReport)
}

// Set godoc
//
// @Router /api/reports/get [get]
// @Summary Получение списка заданий на смену
// @Description При обращении, возвращает список заданий на смену
//
// @Tags Reports
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} []reports "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiReports(w http.ResponseWriter, r *http.Request) {
	reports, err := transport.service.GetReports()
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get reports: %w", err))
		return
	}

	if reports == nil {
		utils.WriteString(w, http.StatusNotFound, "Reports not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, reports)
}

// Set godoc
//
// @Router /api/reports/get/{id} [get]
// @Summary Получение данных о задании на смену
// @Description При обращении, возвращает данные о задании на смену
//
// @Tags Assignments
// @Produce      application/json
// @Consume      application/json
//
// @Params id path int true "ID задания на смену"
//
// @Success 200 {object} reports "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiReportsById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	reports, err := transport.service.GetReportsByReport(idStr)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get reports by id: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, reports)
}

// Set godoc
//
// @Router /api/reports/update [put]
// @Summary Обновление данных о задании на смену
// @Description При обращении, обновляет данные о задании на смену
//
// @Tags Reports
// @Produce      application/json
// @Consume      application/json
//
// @Param request body reports true "Тело запроса"
//
// @Success 200 {object} reports "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PutApiReportsUpdate(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid read body: %w", err))
		return
	}

	var report models.Reports
	if err := jsoniter.Unmarshal(data, &report); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid unmarshal: %w", err))
		return
	}

	ok, err := transport.service.UpdateReport(report)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update report: %w", err))
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update report: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, "OK")
}
