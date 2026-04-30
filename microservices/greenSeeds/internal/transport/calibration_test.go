package transport_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Impisigmatus/service_core/log"
	"github.com/go-chi/chi/v5"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/mocks"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/transport"
	"github.com/rs/zerolog"
)

func TestPostApiCalibrationHandshake(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		mockApp.EXPECT().
			CalibrationHandshake().
			Return("session-123", nil)

		tr := transport.Transport{Calibrate: mockApp}

		req := httptest.NewRequest(http.MethodPost, "/api/calibration/handshake", nil)
		w := httptest.NewRecorder()

		tr.PostApiCalibrationHandshake(w, req)

		res := w.Result()

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		assert.Equal(t, "session-123", res.Header.Get("X-Calibration-Session"))
	})

	t.Run("Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		mockApp.EXPECT().
			CalibrationHandshake().
			Return("", errors.New("fail"))

		tr := transport.Transport{Calibrate: mockApp}

		req := httptest.NewRequest(http.MethodPost, "/api/calibration/handshake", nil)
		w := httptest.NewRecorder()

		tr.PostApiCalibrationHandshake(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("HeaderExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		mockApp.EXPECT().
			CalibrationHandshake().
			Return("abc", nil)

		tr := transport.Transport{Calibrate: mockApp}

		req := httptest.NewRequest(http.MethodPost, "/handshake", nil)
		w := httptest.NewRecorder()

		tr.PostApiCalibrationHandshake(w, req)

		assert.NotEmpty(t, w.Result().Header.Get("X-Calibration-Session"))
	})
}

func TestPostApiCalibrationPhoto(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		mockApp.EXPECT().
			GetPhoto("session-1", "1").
			Return([]byte("image"), nil)

		tr := transport.Transport{Calibrate: mockApp}

		req := httptest.NewRequest(http.MethodPost, "/photo/1", nil)
		req.Header.Set("X-Calibration-Session", "session-1")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("number-of-photo", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		tr.PostApiCalibrationPhoto(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("NoSession", func(t *testing.T) {
		tr := transport.Transport{}

		req := httptest.NewRequest(http.MethodPost, "/photo/1", nil)
		w := httptest.NewRecorder()

		tr.PostApiCalibrationPhoto(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("AppError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		mockApp.EXPECT().
			GetPhoto("session-1", "1").
			Return(nil, errors.New("fail"))

		tr := transport.Transport{Calibrate: mockApp}

		req := httptest.NewRequest(http.MethodPost, "/photo/1", nil)
		req.Header.Set("X-Calibration-Session", "session-1")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("number-of-photo", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		tr.PostApiCalibrationPhoto(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("InvalidPhotoNumber", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		mockApp.EXPECT().
			GetPhoto("session-1", "abc").
			Return(nil, errors.New("bad number"))

		tr := transport.Transport{Calibrate: mockApp}

		req := httptest.NewRequest(http.MethodPost, "/photo/abc", nil)
		req.Header.Set("X-Calibration-Session", "session-1")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("number-of-photo", "abc")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		tr.PostApiCalibrationPhoto(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPostApiCalibrationClear(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		mockApp.EXPECT().Clear("session-1").Return(nil)

		tr := transport.Transport{Calibrate: mockApp}

		req := httptest.NewRequest(http.MethodPost, "/clear", nil)
		req.Header.Set("X-Calibration-Session", "session-1")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationClear(w, req)

		assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
	})

	t.Run("NoSession", func(t *testing.T) {
		tr := transport.Transport{}

		req := httptest.NewRequest(http.MethodPost, "/clear", nil)
		w := httptest.NewRecorder()

		tr.PostApiCalibrationClear(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("AppError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		mockApp.EXPECT().Clear("session-1").Return(errors.New("fail"))

		tr := transport.Transport{Calibrate: mockApp}

		req := httptest.NewRequest(http.MethodPost, "/clear", nil)
		req.Header.Set("X-Calibration-Session", "session-1")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationClear(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("EmptySession", func(t *testing.T) {
		tr := transport.Transport{}

		req := httptest.NewRequest(http.MethodPost, "/clear", nil)
		req.Header.Set("X-Calibration-Session", "")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationClear(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPostApiCalibrationCalc(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		steps := float64(10)

		mockApp.EXPECT().
			CalculateResult(gomock.Any()).
			DoAndReturn(func(c models.Calibration) (models.Calibration, error) {
				assert.Equal(t, "session-1", c.SessionId) // 🔥 важно
				return models.Calibration{
					SessionId: c.SessionId,
					Steps:     &steps,
				}, nil
			})

		tr := transport.Transport{Calibrate: mockApp}

		body := `{
			"steps": 10
		}`

		req := httptest.NewRequest(http.MethodPost, "/calc", bytes.NewBufferString(body))
		req.Header.Set("X-Calibration-Session", "session-1")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationCalc(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("NoSession", func(t *testing.T) {
		tr := transport.Transport{}

		req := httptest.NewRequest(http.MethodPost, "/calc", nil)
		w := httptest.NewRecorder()

		tr.PostApiCalibrationCalc(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		tr := transport.Transport{}

		req := httptest.NewRequest(http.MethodPost, "/calc", bytes.NewBufferString("{bad json"))
		req.Header.Set("X-Calibration-Session", "session-1")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationCalc(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("AppError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		mockApp.EXPECT().
			CalculateResult(gomock.Any()).
			Return(models.Calibration{}, errors.New("fail"))

		tr := transport.Transport{Calibrate: mockApp}

		body := `{"steps": 10}`

		req := httptest.NewRequest(http.MethodPost, "/calc", bytes.NewBufferString(body))
		req.Header.Set("X-Calibration-Session", "session-1")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationCalc(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("EmptyBody", func(t *testing.T) {
		tr := transport.Transport{}

		req := httptest.NewRequest(http.MethodPost, "/calc", nil)
		req.Header.Set("X-Calibration-Session", "session-1")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationCalc(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("SessionInjected", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		mockApp.EXPECT().
			CalculateResult(gomock.Any()).
			DoAndReturn(func(c models.Calibration) (models.Calibration, error) {
				assert.Equal(t, "session-1", c.SessionId)
				return c, nil
			})

		tr := transport.Transport{Calibrate: mockApp}

		body := `{"steps": 5}`

		req := httptest.NewRequest(http.MethodPost, "/calc", bytes.NewBufferString(body))
		req.Header.Set("X-Calibration-Session", "session-1")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationCalc(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
func TestPostApiCalibrationSave(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		mockApp.EXPECT().Save("session-1").Return(nil)

		tr := transport.Transport{Calibrate: mockApp}

		req := httptest.NewRequest(http.MethodPost, "/save", nil)
		req.Header.Set("X-Calibration-Session", "session-1")

		// кладём логгер в контекст
		ctx := context.WithValue(req.Context(), log.CtxKey, zerolog.Nop())
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PostApiCalibrationSave(w, req)

		assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
	})

	t.Run("NoLogger", func(t *testing.T) {
		tr := transport.Transport{}

		req := httptest.NewRequest(http.MethodPost, "/save", nil)
		req.Header.Set("X-Calibration-Session", "session-1")

		w := httptest.NewRecorder()

		tr.PostApiCalibrationSave(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("AppError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		mockApp.EXPECT().Save("session-1").Return(errors.New("fail"))

		tr := transport.Transport{Calibrate: mockApp}

		req := httptest.NewRequest(http.MethodPost, "/save", nil)
		req.Header.Set("X-Calibration-Session", "session-1")

		ctx := context.WithValue(req.Context(), log.CtxKey, zerolog.Nop())
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PostApiCalibrationSave(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("EmptySession", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockICalibrationApp(ctrl)

		tr := transport.Transport{Calibrate: mockApp}

		req := httptest.NewRequest(http.MethodPost, "/save", nil)

		ctx := context.WithValue(req.Context(), log.CtxKey, zerolog.Nop())
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PostApiCalibrationSave(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}