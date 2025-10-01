package transport

// import (
// 	"io"
// 	"net/http"

// 	jsoniter "github.com/json-iterator/go"
// 	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
// 	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/utils"
// )

// // Set godoc
// //
// // @Router /api/register [post]
// // @Summary Регистрация пользователя
// // @Description При обращении, регистрирует пользователя
// //
// // @Tags Auth
// // @Produce      application/json
// // @Consume      application/json
// //
// // @Param user body models.User true "User object"
// //
// // @Success 204 {object} nil "Запрос выполнен успешно"
// // @Failure 400 {object} nil "Ошибка валидации данных"
// // @Failure 401 {object} nil "Ошибка авторизации"
// // @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
// func (transport *Transport) PostApiRegisterUser(w http.ResponseWriter, r *http.Request) {
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	var regUser models.User
// 	if err := jsoniter.Unmarshal(body, &regUser); err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	status, err := transport.service.RegisterUser(regUser)
// 	if err != nil {
// 		utils.WriteString(w, status, err.Error())
// 		return
// 	}

// 	utils.WriteNoContent(w)
// }

// // Set godoc
// //
// // @Router /api/login [post]
// // @Summary Авторизация пользователя
// // @Description При обращении, авторизует пользователя
// //
// // @Tags Auth
// // @Produce      application/json
// // @Consume      application/json
// //
// // @Param user body models.User true "User object"
// //
// // @Success 200 {object} models.TokenResponse "Запрос выполнен успешно"
// // @Failure 400 {object} nil "Ошибка валидации данных"
// // @Failure 401 {object} nil "Ошибка авторизации"
// // @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
// func (transport *Transport) PostApiLoginUser(w http.ResponseWriter, r *http.Request) {
// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	var userLogin models.User
// 	if err := jsoniter.Unmarshal(data, &userLogin); err != nil {
// 		utils.WriteString(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	tokenResponse, status, err := transport.service.LoginUser(userLogin)
// 	if err != nil {
// 		utils.WriteString(w, status, err.Error())
// 		return
// 	}

// 	utils.WriteJSON(w, status, tokenResponse)
// }
