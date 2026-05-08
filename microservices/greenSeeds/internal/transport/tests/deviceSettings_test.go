package transport_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"

	"github.com/Impisigmatus/service_core/log"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/mocks"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/transport"
)

func TestDeviceSettingsTransport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDS := mocks.NewMockIDeviceSettingsApp(ctrl)

	tr := &transport.Transport{
		DeviceSettings: mockDS,
	}

	// ---------- ADD ----------
	t.Run("Add/success", func(t *testing.T) {
		ds := models.DeviceSettings{Key: "k", Value: "v"}
		body, _ := jsoniter.Marshal(ds)

		mockDS.EXPECT().
			AddSetting(gomock.Any()).
			Return(ds, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/device-settings/add", bytes.NewBuffer(body))

		logger := zerolog.Nop()
		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PostApiDeviceSettingsAdd(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Add/no logger", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/device-settings/add", nil)
		w := httptest.NewRecorder()

		tr.PostApiDeviceSettingsAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Add/invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/device-settings/add", bytes.NewBuffer([]byte("bad json")))

		logger := zerolog.Nop()
		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PostApiDeviceSettingsAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Add/app error", func(t *testing.T) {
		ds := models.DeviceSettings{Key: "k", Value: "v"}
		body, _ := jsoniter.Marshal(ds)

		mockDS.EXPECT().
			AddSetting(gomock.Any()).
			Return(models.DeviceSettings{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/api/device-settings/add", bytes.NewBuffer(body))

		logger := zerolog.Nop()
		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PostApiDeviceSettingsAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- GET ALL ----------
	t.Run("Get/success", func(t *testing.T) {
		mockDS.EXPECT().
			GetSettings().
			Return([]models.DeviceSettings{{Key: "k"}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/device-settings/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiDeviceSettingsGet(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Get/not found", func(t *testing.T) {
		mockDS.EXPECT().
			GetSettings().
			Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/device-settings/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiDeviceSettingsGet(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	// ---------- GET BY KEY ----------
	t.Run("GetByKey/success", func(t *testing.T) {
		mockDS.EXPECT().
			GetSettingsByKey("k").
			Return(models.DeviceSettings{Key: "k"}, nil)

		r := chi.NewRouter()
		r.Get("/{key}", tr.GetApiDeviceSettingsGetKey)

		req := httptest.NewRequest(http.MethodGet, "/k", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("GetByKey/error", func(t *testing.T) {
		mockDS.EXPECT().
			GetSettingsByKey("k").
			Return(models.DeviceSettings{}, errors.New("fail"))

		r := chi.NewRouter()
		r.Get("/{key}", tr.GetApiDeviceSettingsGetKey)

		req := httptest.NewRequest(http.MethodGet, "/k", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- UPDATE ----------
	t.Run("Update/success", func(t *testing.T) {
		ds := models.DeviceSettings{Key: "k", Value: "v"}
		body, _ := jsoniter.Marshal(ds)

		mockDS.EXPECT().
			UpdateSetting(gomock.Any()).
			Return(ds, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/device-settings/update", bytes.NewBuffer(body))

		logger := zerolog.Nop()
		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PutApiDeviceSettingsUpdate(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Update/no logger", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/device-settings/update", nil)
		w := httptest.NewRecorder()

		tr.PutApiDeviceSettingsUpdate(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Update/app error", func(t *testing.T) {
		ds := models.DeviceSettings{Key: "k"}
		body, _ := jsoniter.Marshal(ds)

		mockDS.EXPECT().
			UpdateSetting(gomock.Any()).
			Return(models.DeviceSettings{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPut, "/api/device-settings/update", bytes.NewBuffer(body))

		logger := zerolog.Nop()
		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PutApiDeviceSettingsUpdate(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- DELETE ----------
	t.Run("Delete/success", func(t *testing.T) {
		mockDS.EXPECT().
			DeleteSetting("k").
			Return(true, nil)

		r := chi.NewRouter()
		r.Delete("/{key}", tr.DeleteApiDeviceSettingsDelete)

		req := httptest.NewRequest(http.MethodDelete, "/k", nil)

		logger := zerolog.Nop()
		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}
	})

	t.Run("Delete/no logger", func(t *testing.T) {
		r := chi.NewRouter()
		r.Delete("/{key}", tr.DeleteApiDeviceSettingsDelete)

		req := httptest.NewRequest(http.MethodDelete, "/k", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Delete/fail", func(t *testing.T) {
		mockDS.EXPECT().
			DeleteSetting("k").
			Return(false, nil)

		r := chi.NewRouter()
		r.Delete("/{key}", tr.DeleteApiDeviceSettingsDelete)

		req := httptest.NewRequest(http.MethodDelete, "/k", nil)

		logger := zerolog.Nop()
		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}
