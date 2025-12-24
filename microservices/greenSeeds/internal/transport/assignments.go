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
// @Router /api/assignments/add [post]
// @Summary Добавление информации о задании на смену
// @Description При обращении, добавляет информацию о задании на смену в БД
//
// @Tags Assignments
// @Produce      application/json
// @Consume      application/json
//
// @Param 	request	body	assignment	true	"Тело запроса"
//
// @Success 200 {object} assignment "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiAssignmentsAdd(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %w", err))
		return
	}

	var assignment models.Assignments
	if err := jsoniter.Unmarshal(body, &assignment); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %w", err))
		return
	}

	addedAssignment, err := transport.service.AddAssignment(assignment)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add assignment: %w", err))
		return
	}

	if addedAssignment == (models.Assignments{}) {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add assignment: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, addedAssignment)
}

// Set godoc
//
// @Router /api/assignments/get [get]
// @Summary Получение списка заданий на смену
// @Description При обращении, возвращает список заданий на смену
//
// @Tags Assignments
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} []assignment "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiAssignmentsGet(w http.ResponseWriter, r *http.Request) {
	assignments, err := transport.service.GetAssignments()
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get assignments: %w", err))
		return
	}

	if assignments == nil {
		utils.WriteString(w, http.StatusNotFound, "Assignments not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, assignments)
}

// Set godoc
//
// @Router /api/assignments/get/{id} [get]
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
func (transport *Transport) GetApiAssignmentsGetAssignment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	assignment, err := transport.service.GetAssignmentsByAssignment(idStr)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get assignment by assignment: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, assignment)
}

// Set godoc
//
// @Router /api/assignments/update [put]
// @Summary Обновление данных о задании на смену
// @Description При обращении, обновляет данные о задании на смену
//
// @Tags Assignments
// @Produce      application/json
// @Consume      application/json
//
// @Param request body assignment true "Тело запроса"
//
// @Success 200 {object} assignment "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PutApiAssignmentsUpdate(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid read body: %w", err))
		return
	}

	var assignment models.Assignments
	if err := jsoniter.Unmarshal(data, &assignment); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid unmarshal: %w", err))
		return
	}

	updatedAssignment, err := transport.service.UpdateAssignment(assignment)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update assignment: %w", err))
		return
	}

	if updatedAssignment == (models.Assignments{}) {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update assignment: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedAssignment)
}

// Set godoc
//
// @Router /api/assignments/delete/{id} [delete]
// @Summary Удаление задания на смену
// @Description При обращении, удаляет задание на смену
//
// @Tags Assignments
// @Produce      application/json
// @Consume      application/json
//
// @Param id path int true "ID задания на смену"
//
// @Success 204 {object} nil "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) DeleteApiAssignmentsDelete(w http.ResponseWriter, r *http.Request) {
	log, ok := r.Context().Value(log.CtxKey).(zerolog.Logger)
	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid logger")
		return
	}

	idStr := chi.URLParam(r, "id")

	ok, err := transport.service.DeleteAssignments(idStr)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid delete assignment")
		return
	}

	log.Warn().Ctx(r.Context()).Msg(fmt.Sprintf("Assignment removed: %s", idStr))

	utils.WriteNoContent(w)
}

// Set godoc
//
// @Router /api/assignments/active-tasks/{username} [get]
// @Summary Получение списка активных заданий
// @Description При обращении, возвращает список активных заданий
//
// @Tags Assignments
// @Produce      application/json
// @Consume      application/json
//
// @Param username path string true "Имя пользователя"
//
// @Success 200 {object} []active_task "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiActiveTasks(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	check, err := transport.service.CheckActiveTasks(username)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get active task: %w", err))
		return
	}

	if check == nil {
		utils.WriteString(w, http.StatusNotFound, "Active tasks not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, check)
}

// Set godoc
//
// @Router /api/assignments/task/{id} [get]
// @Summary Получение списка активных заданий
// @Description При обращении, возвращает список активных заданий
//
// @Tags Assignments
// @Produce      application/json
// @Consume      application/json
//
// @Param id path int true "ID задания"
//
// @Success 200 {object} []active_task "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	task, err := transport.service.GetTaskById(idStr)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get task: %w", err))
		return
	}

	if task == (models.Task{}) {
		utils.WriteString(w, http.StatusNotFound, "Task not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, task)
}
