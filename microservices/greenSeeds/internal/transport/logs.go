package transport

import (
	"fmt"
	"net/http"
	"time"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/utils"
)

const layout = "2006-01-02"

// Set godoc
//
// @Router /api/logs/get [get]
// @Summary Получение логов из БД
// @Description При обращении, возвращает список заданий на смену
//
// @Tags Logs
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} []assignment "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiLogsGet(w http.ResponseWriter, r *http.Request) {
	var params models.LogsParams

	search := r.URL.Query().Get("search")
	if search != "" {
		params.Search = fmt.Sprintf("%%%s%%", search)
	}

	params.Level = r.URL.Query().Get("level")
	params.Limit = r.URL.Query().Get("limit")
	params.Offset = r.URL.Query().Get("offset")

	dateFrom := r.URL.Query().Get("date_from")
	if dateFrom == "" {
		params.DateFrom = nil
	} else {
		timeFrom, err := time.Parse(layout, dateFrom)
		if err != nil {
			params.DateFrom = nil
		}
		params.DateFrom = &timeFrom
	}

	dateTo := r.URL.Query().Get("date_to")
	if dateTo == "" {
		params.DateTo = nil
	} else {
		timeTo, err := time.Parse(layout, dateTo)
		if err != nil {
			params.DateTo = nil
		}
		params.DateTo = &timeTo
	}

	logs, err := transport.service.GetLogs(params)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get logs: %w", err))
		return
	}

	if logs == nil {
		utils.WriteString(w, http.StatusNotFound, "Logs not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, logs)
}
