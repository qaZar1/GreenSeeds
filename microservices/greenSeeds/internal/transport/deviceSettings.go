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
// @Router /api/device-settings/add [post]
// @Summary Добавление настроек в БД
// @Description При обращении, добавляет настройки в БД
//
// @Tags Bunkers
// @Produce      application/json
// @Consume      application/json
//
// @Param 	request	body	deviceSettings	true	"Тело запроса"
//
// @Success 200 {object} deviceSettings "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PostApiDeviceSettingsAdd(w http.ResponseWriter, r *http.Request) {
	log, ok := r.Context().Value(log.CtxKey).(zerolog.Logger)
	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid logger")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not read body: %w", err))
		return
	}

	var deviceSettings models.DeviceSettings
	if err := jsoniter.Unmarshal(body, &deviceSettings); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Can not unmarshal: %w", err))
		return
	}

	inserted, err := transport.service.AddSetting(deviceSettings)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add device settings: %w", err))
		return
	}

	if inserted == (models.DeviceSettings{}) {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid add device settings: %w", err))
		return
	}

	log.Warn().Ctx(r.Context()).Msg(fmt.Sprintf("Device settings added %s: %s", inserted.Key, inserted.Value))

	utils.WriteJSON(w, http.StatusOK, inserted)
}

// Set godoc
//
// @Router /api/device-settings/get [get]
// @Summary Получение списка настроек
// @Description При обращении, возвращает список настроек
//
// @Tags DeviceSettings
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} []deviceSettings "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiDeviceSettingsGet(w http.ResponseWriter, r *http.Request) {
	deviceSettings, err := transport.service.GetSettings()
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get device settings: %w", err))
		return
	}

	if deviceSettings == nil {
		utils.WriteString(w, http.StatusNotFound, fmt.Sprintf("Device settings not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, deviceSettings)
}

// Set godoc
//
// @Router /api/device-settings/get/{key} [get]
// @Summary Получение настроек по ключу
// @Description При обращении, возвращает настройки по ключу
//
// @Tags DeviceSettings
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} deviceSettings "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) GetApiDeviceSettingsGetKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	deviceSettings, err := transport.service.GetSettingsByKey(key)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid get device settings by key: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, deviceSettings)
}

// Set godoc
//
// @Router /api/device-settings/update [put]
// @Summary Обновление данных о настройках
// @Description При обращении, обновляет данные о настройках
//
// @Tags DeviceSettings
// @Produce      application/json
// @Consume      application/json
//
// @Success 200 {object} deviceSettings "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) PutApiDeviceSettingsUpdate(w http.ResponseWriter, r *http.Request) {
	log, ok := r.Context().Value(log.CtxKey).(zerolog.Logger)
	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid logger")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid read body: %w", err))
		return
	}

	var deviceSettings models.DeviceSettings
	if err := jsoniter.Unmarshal(data, &deviceSettings); err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid unmarshal: %w", err))
		return
	}

	updated, err := transport.service.UpdateSetting(deviceSettings)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update settings: %w", err))
		return
	}

	if updated == (models.DeviceSettings{}) {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid update settings: %w", err))
		return
	}

	log.Warn().Ctx(r.Context()).Msg(fmt.Sprintf("Device settings updated %s: %s", updated.Key, updated.Value))

	utils.WriteJSON(w, http.StatusOK, updated)
}

// Set godoc
//
// @Router /api/device-settings/delete/{key} [delete]
// @Summary Удаление настроек
// @Description При обращении, удаляет настройки
//
// @Tags DeviceSettings
// @Produce      application/json
// @Consume      application/json
//
// @Param key path string true "Ключ настроек"
//
// @Success 204 {object} nil "Запрос выполнен успешно"
// @Failure 400 {object} nil "Ошибка валидации данных"
// @Failure 401 {object} nil "Ошибка авторизации"
// @Failure 500 {object} nil "Произошла внутренняя ошибка сервера"
func (transport *Transport) DeleteApiDeviceSettingsDelete(w http.ResponseWriter, r *http.Request) {
	log, ok := r.Context().Value(log.CtxKey).(zerolog.Logger)
	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, "Invalid logger")
		return
	}

	key := chi.URLParam(r, "key")

	ok, err := transport.service.DeleteSetting(key)
	if err != nil {
		utils.WriteString(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !ok {
		utils.WriteString(w, http.StatusInternalServerError, fmt.Sprintf("Invalid delete settings: %w", err))
		return
	}

	log.Warn().Ctx(r.Context()).Msg(fmt.Sprintf("Device settings deleted %s", key))

	utils.WriteNoContent(w)
}
