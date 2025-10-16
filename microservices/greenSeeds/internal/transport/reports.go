package transport

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/utils"
)

// // Set godoc
// //
// // @Router /api/assignments/add [post]
// // @Summary Добавление информации о задании на смену
// // @Description При обращении, добавляет информацию о задании на смену в БД
// //
// // @Tags Assignments
// // @Produce      application/json
// // @Consume      application/json
// //
// // @Param 	request	body	assignment	true	"Тело запроса"
// //
// // @Success 200 {object} assignment "Запрос выполнен успешно"
// // @Failure 400 {object} nil "Ошибка валидации данных"
// // @Failure 401 {object} nil "Ошибка авторизации"
// // @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
// func (transport *Transport) PostApiAssignmentsAdd(w http.ResponseWriter, r *http.Request) {
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %w", err))
// 		return
// 	}

// 	var assignment models.Assignments
// 	if err := jsoniter.Unmarshal(body, &assignment); err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %w", err))
// 		return
// 	}

// 	addedAssignment, err := transport.service.AddAssignment(assignment)
// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add assignment: %w", err))
// 		return
// 	}

// 	if addedAssignment == (models.Assignments{}) {
// 		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add assignment: %w", err))
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, addedAssignment)
// }

// Set godoc
//
// @Router /api/reports/get [get]
// @Summary Получение списка заданий на смену
// @Description При обращении, возвращает список заданий на смену
//
// @Tags Assignments
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

// // Set godoc
// //
// // @Router /api/assignments/update [put]
// // @Summary Обновление данных о задании на смену
// // @Description При обращении, обновляет данные о задании на смену
// //
// // @Tags Assignments
// // @Produce      application/json
// // @Consume      application/json
// //
// // @Param request body assignment true "Тело запроса"
// //
// // @Success 200 {object} assignment "Запрос выполнен успешно"
// // @Failure 400 {object} nil "Ошибка валидации данных"
// // @Failure 401 {object} nil "Ошибка авторизации"
// // @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
// func (transport *Transport) PutApiAssignmentsUpdate(w http.ResponseWriter, r *http.Request) {
// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid read body: %w", err))
// 		return
// 	}

// 	var assignment models.Assignments
// 	if err := jsoniter.Unmarshal(data, &assignment); err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid unmarshal: %w", err))
// 		return
// 	}

// 	updatedAssignment, err := transport.service.UpdateAssignment(assignment)
// 	if err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update assignment: %w", err))
// 		return
// 	}

// 	if updatedAssignment == (models.Assignments{}) {
// 		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update assignment: %w", err))
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, updatedAssignment)
// }

// // Set godoc
// //
// // @Router /api/assignments/delete/{id} [delete]
// // @Summary Удаление задания на смену
// // @Description При обращении, удаляет задание на смену
// //
// // @Tags Assignments
// // @Produce      application/json
// // @Consume      application/json
// //
// // @Param id path int true "ID задания на смену"
// //
// // @Success 204 {object} nil "Запрос выполнен успешно"
// // @Failure 400 {object} nil "Ошибка валидации данных"
// // @Failure 401 {object} nil "Ошибка авторизации"
// // @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
// func (transport *Transport) DeleteApiAssignmentsDelete(w http.ResponseWriter, r *http.Request) {
// 	idStr := chi.URLParam(r, "id")

// 	ok, err := transport.service.DeleteAssignments(idStr)
// 	if err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	if !ok {
// 		utils.WriteString(w, http.StatusInternalServerError, "Invalid delete assignment")
// 		return
// 	}

// 	utils.WriteNoContent(w)
// }
