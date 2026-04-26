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
// @Router /api/users/get/{id} [get]
// @Summary Получение пользователя
// @Description При обращении, возвращает пользователя по его username
//
// @Tags Users
// @Produce      application/json
//
// @Param id path int true "User id"
//
// @Success 200 {object} user "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiUserGetUsername(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user_id")
	user, err := transport.app.GetUserById(userId)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not get user by id: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

// Set godoc
//
// @Router /api/users/get [get]
// @Summary Получение всех пользователей
// @Description При обращении, возвращает всех пользователей
//
// @Tags Users
// @Produce      application/json
//
// @Success 200 {array} user "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiCheckAllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := transport.app.CheckAllUsers()
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not get all users: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, allUsers)
}

// Set godoc
//
// @Router /api/users/change-password [put]
// @Summary Обновление пароля
// @Description При обращении, обновляет пароль
//
// @Tags Users
// @Produce      application/json
// @Consume      application/json
//
// @Param 	request	body	UpdatePassword	true	"Тело запроса"
//
// @Success 204 {object} nil "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PutApiChangePassword(w http.ResponseWriter, r *http.Request) {
	log, ok := r.Context().Value(log.CtxKey).(zerolog.Logger)
	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid logger")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %w", err))
		return
	}

	var updatePassword models.UpdatePassword
	if err := jsoniter.Unmarshal(data, &updatePassword); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %w", err))
		return
	}

	ok, err = transport.app.ChangePassword(updatePassword)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not change password: %w", err))
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusNotFound, "User not found")
		return
	}

	log.Warn().Ctx(r.Context()).Msg(
		fmt.Sprintf(
			"The password has been reset for: %d",
			updatePassword.Id,
		),
	)

	utils.WriteNoContent(w)
}

// Set godoc
//
// @Router /api/users/removeUser [delete]
// @Summary Удаление пользователя
// @Description При обращении, удаляет пользователя
//
// @Tags Users
// @Produce      application/json
// @Consume      application/json
//
// @Param removeUser path string true "User username"
//
// @Success 204 {object} nil "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) DeleteApiRemoveUser(w http.ResponseWriter, r *http.Request) {
	log, ok := r.Context().Value(log.CtxKey).(zerolog.Logger)
	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid logger")
		return
	}

	username := chi.URLParam(r, "username")

	ok, err := transport.app.RemoveUser(username)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusNotFound, "User not found")
		return
	}

	log.Warn().Ctx(r.Context()).Msg(fmt.Sprintf("User removed: %s", username))

	utils.WriteNoContent(w)
}

// Set godoc
//
// @Router /api/users/update [put]
// @Summary Обновление пользователя
// @Description При обращении, обновляет пользователя
//
// @Tags Users
// @Produce      application/json
// @Consume      application/json
//
// @Param 	request	body	user	true	"Тело запроса"
//
// @Success 204 {object} nil "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PutApiUpdateUser(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %w", err))
		return
	}

	var user models.User
	if err := jsoniter.Unmarshal(data, &user); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %w", err))
		return
	}

	ok, err := transport.app.Update(user)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not update user: %w", err))
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusNotFound, "User not found")
		return
	}

	utils.WriteNoContent(w)
}
