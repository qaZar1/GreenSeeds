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
// @Router /api/receipts/add [post]
// @Summary Добавление информации, в каком бункере какие семена
// @Description При обращении, добавляет информацию о семенах в БД
//
// @Tags Receipts
// @Produce      application/json
// @Consume      application/json
//
// @Param 	request	body	receipts	true	"Тело запроса"
//
// @Success 200 {object} receipts "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiReceiptsAdd(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %w", err))
		return
	}

	var receipt models.Receipts
	if err := jsoniter.Unmarshal(body, &receipt); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %w", err))
		return
	}

	addedReceipt, err := transport.service.AddReceipts(receipt)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add receipts: %w", err))
		return
	}

	if addedReceipt == (models.Receipts{}) {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add placement: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, addedReceipt)
}

// Set godoc
//
// @Router /api/receipts/get [get]
// @Summary Получение списка бункеров и семян
// @Description При обращении, возвращает список бункеров и семян
//
// @Tags Receipts
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} []receipts "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiReceiptsGet(w http.ResponseWriter, r *http.Request) {
	receipts, err := transport.service.GetReceipts()
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get receipts: %w", err))
		return
	}

	if receipts == nil {
		utils.WriteString(w, http.StatusNotFound, fmt.Sprintf("Receipts not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, receipts)
}

// Set godoc
//
// @Router /api/receipts/get/{receipt} [get]
// @Summary Получение бункера и семян по ID
// @Description При обращении, возвращает бункер и семена по ID
//
// @Tags Receipts
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} receipts "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiReceiptsGetReceipt(w http.ResponseWriter, r *http.Request) {
	receiptName := chi.URLParam(r, "receipt")

	receipt, err := transport.service.GetReceiptsByReceipt(receiptName)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get receipt by receipt: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, receipt)
}

// Set godoc
//
// @Router /api/receipts/update [put]
// @Summary Обновление данных о семенах
// @Description При обращении, обновляет данные о семенах
//
// @Tags Receipts
// @Produce      application/json
// @Consume      application/json
//
// @Param request body receipts true "Тело запроса"
//
// @Success 200 {object} receipts "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PutApiReceiptsUpdate(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid read body: %w", err))
		return
	}

	var receipt models.Receipts
	if err := jsoniter.Unmarshal(data, &receipt); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid unmarshal: %w", err))
		return
	}

	updatedReceipt, err := transport.service.UpdateReceipts(receipt)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update receipt: %w", err))
		return
	}

	if updatedReceipt == (models.Receipts{}) {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update receipt: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedReceipt)
}

// Set godoc
//
// @Router /api/receipts/delete/{receipt} [delete]
// @Summary Удаление рецепта
// @Description При обращении, удаляет рецепт
//
// @Tags Receipts
// @Produce      application/json
// @Consume      application/json
//
// @Param receipt path string true "Название рецепта"
//
// @Success 204 {object} nil "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) DeleteApiReceiptsDelete(w http.ResponseWriter, r *http.Request) {
	log, ok := r.Context().Value(log.CtxKey).(zerolog.Logger)
	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid logger")
		return
	}

	receiptName := chi.URLParam(r, "receipt")

	ok, err := transport.service.DeleteReceipts(receiptName)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid delete receipt")
		return
	}

	log.Warn().Ctx(r.Context()).Msg(fmt.Sprintf("Receipt removed: %s", receiptName))

	utils.WriteNoContent(w)
}
