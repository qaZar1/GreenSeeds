package transport

// import (
// 	"io"
// 	"net/http"

// 	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/utils"
// 	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/vendor/github.com/go-chi/chi/v5"
// 	jsoniter "github.com/qaZar1/GreenSeeds/microservices/greenSeeds/vendor/github.com/json-iterator/go"
// 	"github.com/qaZar1/ITops/microservices/itops/internal/models"
// )

// // import (
// // 	"io"
// // 	"net/http"

// // 	"github.com/go-chi/chi/v5"
// // 	jsoniter "github.com/json-iterator/go"
// // 	"github.com/qaZar1/ITops/microservices/itops/internal/models"
// // 	"github.com/qaZar1/ITops/microservices/itops/internal/utils"
// // )

// // Set godoc
// //
// // @Router /api/users/checkByUuid/{uuid} [get]
// // @Summary Получение пользователя
// // @Description При обращении, возвращает пользователя по его uuid
// //
// // @Tags Users
// // @Produce      application/json
// //
// // @Param uuid path string true "User uuid"
// //
// // @Success 200 {object} models.User "Запрос выполнен успешно"
// // @Failure 400 {object} nil "Ошибка валидации данных"
// // @Failure 401 {object} nil "Ошибка авторизации"
// // @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
// func (transport *Transport) GetApiCheckUserByUuidUuid(w http.ResponseWriter, r *http.Request) {
// 	uuid := chi.URLParam(r, "uuid")
// 	user, status, err := transport.service.CheckUserByUuid(uuid)
// 	if err != nil {
// 		utils.WriteString(w, status, err.Error())
// 		return
// 	}

// 	utils.WriteJSON(w, status, user)
// }

// // Set godoc
// //
// // @Router /api/users/checkRoles/{uuid} [get]
// // @Summary Получение ролей
// // @Description При обращении, возвращает роли по uuid
// //
// // @Tags Users
// // @Produce      application/json
// //
// // @Param uuid path string true "User uuid"
// //
// // @Success 200 {object} models.User "Запрос выполнен успешно"
// // @Failure 400 {object} nil "Ошибка валидации данных"
// // @Failure 401 {object} nil "Ошибка авторизации"
// // @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
// func (transport *Transport) GetApiCheckRolesUuid(w http.ResponseWriter, r *http.Request) {
// 	uuid := chi.URLParam(r, "uuid")
// 	roles, status, err := transport.service.CheckRolesById(uuid)
// 	if err != nil {
// 		utils.WriteString(w, status, err.Error())
// 		return
// 	}

// 	utils.WriteJSON(w, status, roles)
// }

// // Set godoc
// //
// // @Router /api/users/checkAll [get]
// // @Summary Получение всех пользователей
// // @Description При обращении, возвращает всех пользователей
// //
// // @Tags Users
// // @Produce      application/json
// //
// // @Success 200 {array} models.User "Запрос выполнен успешно"
// // @Failure 400 {object} nil "Ошибка валидации данных"
// // @Failure 401 {object} nil "Ошибка авторизации"
// // @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
// func (transport *Transport) GetApiCheckAllUsers(w http.ResponseWriter, r *http.Request) {
// 	allUsers, status, err := transport.service.CheckAllUsers()
// 	if err != nil {
// 		utils.WriteString(w, status, err.Error())
// 		return
// 	}

// 	utils.WriteJSON(w, status, allUsers)
// }

// // Set godoc
// //
// // @Router /api/users/updateRole [put]
// // @Summary Обновление роли
// // @Description При обращении, обновляет роль
// //
// // @Tags Users
// // @Produce      application/json
// // @Consume      application/json
// //
// // @Param updateRole body models.UpdateRole true "Update role object"
// //
// // @Success 204 {object} models.UpdateRole "Запрос выполнен успешно"
// // @Failure 400 {object} nil "Ошибка валидации данных"
// // @Failure 401 {object} nil "Ошибка авторизации"
// // @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
// func (transport *Transport) PutApiUpdateRole(w http.ResponseWriter, r *http.Request) {
// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	var updateRole models.UpdateRole
// 	if err := jsoniter.Unmarshal(data, &updateRole); err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	result, status, err := transport.service.UpdateRole(updateRole)
// 	if err != nil {
// 		utils.WriteString(w, status, err.Error())
// 		return
// 	}

// 	utils.WriteJSON(w, status, result)
// }

// // Set godoc
// //
// // @Router /api/users/change-password [put]
// // @Summary Обновление пароля
// // @Description При обращении, обновляет пароль
// //
// // @Tags Users
// // @Produce      application/json
// // @Consume      application/json
// //
// // @Param updatePassword body models.UpdatePassword true "Update password object"
// //
// // @Success 204 {object} nil "Запрос выполнен успешно"
// // @Failure 400 {object} nil "Ошибка валидации данных"
// // @Failure 401 {object} nil "Ошибка авторизации"
// // @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
// func (transport *Transport) PutApiChangePassword(w http.ResponseWriter, r *http.Request) {
// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	var updatePassword models.UpdatePassword
// 	if err := jsoniter.Unmarshal(data, &updatePassword); err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	status, err := transport.service.ChangePassword(updatePassword)
// 	if err != nil {
// 		utils.WriteString(w, status, err.Error())
// 		return
// 	}

// 	utils.WriteNoContent(w)
// }

// // Set godoc
// //
// // @Router /api/users/reset-password/{uuid} [put]
// // @Summary Обновление пароля
// // @Description При обращении, обновляет пароль
// //
// // @Tags Users
// // @Produce      application/json
// // @Consume      application/json
// //
// // @Param uuid path string true "User uuid"
// //
// // @Success 204 {object} nil "Запрос выполнен успешно"
// // @Failure 400 {object} nil "Ошибка валидации данных"
// // @Failure 401 {object} nil "Ошибка авторизации"
// // @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
// func (transport *Transport) PutApiResetPassword(w http.ResponseWriter, r *http.Request) {
// 	uuid := chi.URLParam(r, "uuid")

// 	if uuid == "" {
// 		utils.WriteString(w, http.StatusBadRequest, ErrInvalidGeneratePassHash.Error())
// 		return
// 	}

// 	status, err := transport.service.ResetPassword(uuid)
// 	if err != nil {
// 		utils.WriteString(w, status, err.Error())
// 		return
// 	}

// 	utils.WriteNoContent(w)
// }

// // Set godoc
// //
// // @Router /api/users/removeUser [delete]
// // @Summary Удаление пользователя
// // @Description При обращении, удаляет пользователя
// //
// // @Tags Users
// // @Produce      application/json
// // @Consume      application/json
// //
// // @Param removeUser body models.RemoveUser true "Remove user object"
// //
// // @Success 204 {object} nil "Запрос выполнен успешно"
// // @Failure 400 {object} nil "Ошибка валидации данных"
// // @Failure 401 {object} nil "Ошибка авторизации"
// // @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
// func (transport *Transport) DeleteApiRemoveUser(w http.ResponseWriter, r *http.Request) {
// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	var removeUser models.RemoveUser
// 	if err := jsoniter.Unmarshal(data, &removeUser); err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	status, err := transport.service.RemoveUser(removeUser)
// 	if err != nil {
// 		utils.WriteString(w, status, err.Error())
// 		return
// 	}

// 	utils.WriteNoContent(w)
// }
