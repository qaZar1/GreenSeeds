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

func TestCalibrationTransport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCal := mocks.NewMockICalibrationApp(ctrl)

	tr := &transport.Transport{
		Calibration: mockCal,
	}

	// ---------- HANDSHAKE ----------
	t.Run("Handshake/success", func(t *testing.T) {
		mockCal.EXPECT().
			CalibrationHandshake().
			Return("session123", nil)

		req := httptest.NewRequest(http.MethodPost, "/api/calibration/handshake", nil)
		w := httptest.NewRecorder()

		tr.PostApiCalibrationHandshake(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}

		if w.Header().Get("X-Calibration-Session") != "session123" {
			t.Fatalf("header not set")
		}
	})

	t.Run("Handshake/error", func(t *testing.T) {
		mockCal.EXPECT().
			CalibrationHandshake().
			Return("", errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/api/calibration/handshake", nil)
		w := httptest.NewRecorder()

		tr.PostApiCalibrationHandshake(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- PHOTO ----------
	t.Run("Photo/success", func(t *testing.T) {
		mockCal.EXPECT().
			GetPhoto("session123", "1").
			Return([]byte("img"), nil)

		r := chi.NewRouter()
		r.Post("/{number-of-photo}", tr.PostApiCalibrationPhoto)

		req := httptest.NewRequest(http.MethodPost, "/1", nil)
		req.Header.Set("X-Calibration-Session", "session123")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Photo/no session", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/calibration/photo/1", nil)
		w := httptest.NewRecorder()

		tr.PostApiCalibrationPhoto(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Photo/error", func(t *testing.T) {
		mockCal.EXPECT().
			GetPhoto("session123", "1").
			Return(nil, errors.New("fail"))

		r := chi.NewRouter()
		r.Post("/{number-of-photo}", tr.PostApiCalibrationPhoto)

		req := httptest.NewRequest(http.MethodPost, "/1", nil)
		req.Header.Set("X-Calibration-Session", "session123")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- CLEAR ----------
	t.Run("Clear/success", func(t *testing.T) {
		mockCal.EXPECT().
			Clear("session123").
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/api/calibration/clear", nil)
		req.Header.Set("X-Calibration-Session", "session123")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationClear(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}
	})

	t.Run("Clear/error", func(t *testing.T) {
		mockCal.EXPECT().
			Clear("session123").
			Return(errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/api/calibration/clear", nil)
		req.Header.Set("X-Calibration-Session", "session123")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationClear(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- CALC ----------
	t.Run("Calc/success", func(t *testing.T) {
		body, _ := jsoniter.Marshal(models.Calibration{})

		mockCal.EXPECT().
			CalculateResult(gomock.Any()).
			Return(models.Calibration{SessionId: "session123"}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/calibration/calc", bytes.NewBuffer(body))
		req.Header.Set("X-Calibration-Session", "session123")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationCalc(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Calc/no session", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/calibration/calc", nil)
		w := httptest.NewRecorder()

		tr.PostApiCalibrationCalc(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Calc/bad json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/calibration/calc", bytes.NewBuffer([]byte("bad")))
		req.Header.Set("X-Calibration-Session", "session123")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationCalc(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- SAVE ----------
	t.Run("Save/success", func(t *testing.T) {
		mockCal.EXPECT().
			Save("session123").
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/api/calibration/save", nil)
		req.Header.Set("X-Calibration-Session", "session123")

		logger := zerolog.Nop()
		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PostApiCalibrationSave(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}
	})

	t.Run("Save/no logger", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/calibration/save", nil)
		req.Header.Set("X-Calibration-Session", "session123")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationSave(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Save/error", func(t *testing.T) {
		mockCal.EXPECT().
			Save("session123").
			Return(errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/api/calibration/save", nil)
		req.Header.Set("X-Calibration-Session", "session123")

		logger := zerolog.Nop()
		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PostApiCalibrationSave(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}
