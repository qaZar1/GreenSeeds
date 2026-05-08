package transport_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"

	"github.com/Impisigmatus/service_core/log"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/mocks"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/transport"
	"github.com/rs/zerolog"
)

func TestShiftsTransport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShifts := mocks.NewMockIShiftsApp(ctrl)

	tr := &transport.Transport{
		Shifts: mockShifts,
	}

	now := time.Now()
	shiftID := int64(1)
	username := "user"
	userID := int64(10)

	// ---------- ADD ----------
	t.Run("Add/success", func(t *testing.T) {
		s := models.Shifts{
			Shift:    &shiftID,
			Dt:       now,
			Username: &username,
			UserId:   &userID,
		}
		body, _ := jsoniter.Marshal(s)

		mockShifts.EXPECT().
			AddShift(gomock.Any()).
			Return(s, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/shifts/add", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiShiftAdd(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Add/error", func(t *testing.T) {
		s := models.Shifts{Dt: now}
		body, _ := jsoniter.Marshal(s)

		mockShifts.EXPECT().
			AddShift(gomock.Any()).
			Return(models.Shifts{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/api/shifts/add", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiShiftAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- GET ----------
	t.Run("Get/success", func(t *testing.T) {
		mockShifts.EXPECT().
			GetShifts().
			Return([]models.Shifts{{Dt: now}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/shifts/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiShiftsGet(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Get/not found", func(t *testing.T) {
		mockShifts.EXPECT().
			GetShifts().
			Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/shifts/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiShiftsGet(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	// ---------- GET BY SHIFT ----------
	t.Run("GetByShift/success", func(t *testing.T) {
		mockShifts.EXPECT().
			GetShiftsByShift("1").
			Return(models.Shifts{Dt: now}, nil)

		r := chi.NewRouter()
		r.Get("/{shift}", tr.GetApiShiftsGetShift)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	// ---------- UPDATE ----------
	t.Run("Update/success", func(t *testing.T) {
		s := models.Shifts{Dt: now}
		body, _ := jsoniter.Marshal(s)

		mockShifts.EXPECT().
			UpdateShifts(gomock.Any()).
			Return(s, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/shifts/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiShiftsUpdate(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Update/error", func(t *testing.T) {
		s := models.Shifts{Dt: now}
		body, _ := jsoniter.Marshal(s)

		mockShifts.EXPECT().
			UpdateShifts(gomock.Any()).
			Return(models.Shifts{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPut, "/api/shifts/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiShiftsUpdate(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- DELETE ----------
	t.Run("Delete/success", func(t *testing.T) {
		mockShifts.EXPECT().
			DeleteShifts("1").
			Return(true, nil)

		r := chi.NewRouter()
		logger := zerolog.Nop()

		r.With(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), log.CtxKey, logger)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}).Delete("/{shift}", tr.DeleteApiShiftsDelete)

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}
	})

	t.Run("Delete/fail", func(t *testing.T) {
		mockShifts.EXPECT().
			DeleteShifts("1").
			Return(false, nil)

		r := chi.NewRouter()
		logger := zerolog.Nop()

		r.With(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), log.CtxKey, logger)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}).Delete("/{shift}", tr.DeleteApiShiftsDelete)

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- GET WITHOUT USER ----------
	t.Run("GetWithoutUser/success", func(t *testing.T) {
		mockShifts.EXPECT().
			GetShiftsWithoutUser().
			Return([]models.Shifts{{Dt: now}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/shifts/getWithoutUser", nil)
		w := httptest.NewRecorder()

		tr.GetApiShiftsGetWithoutUser(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("GetWithoutUser/not found", func(t *testing.T) {
		mockShifts.EXPECT().
			GetShiftsWithoutUser().
			Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/shifts/getWithoutUser", nil)
		w := httptest.NewRecorder()

		tr.GetApiShiftsGetWithoutUser(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}
