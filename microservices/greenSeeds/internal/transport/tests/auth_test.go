package transport_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"

	"github.com/Impisigmatus/service_core/log"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/mocks"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/transport"
)

func TestAuthTransport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsers := mocks.NewMockIUsersApp(ctrl)

	tr := &transport.Transport{
		Users: mockUsers,
	}

	t.Run("Register/success", func(t *testing.T) {
		user := models.User{Username: "test"}
		body, _ := jsoniter.Marshal(user)

		mockUsers.EXPECT().
			RegisterUser(gomock.Any()).
			Return(http.StatusNoContent, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))

		logger := zerolog.Nop()
		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PostApiRegisterUser(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}
	})

	t.Run("Register/no logger", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/register", nil)
		w := httptest.NewRecorder()

		tr.PostApiRegisterUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Register/bad body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/register", nil)

		logger := zerolog.Nop()
		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PostApiRegisterUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Register/invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer([]byte("bad json")))

		logger := zerolog.Nop()
		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PostApiRegisterUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Register/app error", func(t *testing.T) {
		user := models.User{Username: "test"}
		body, _ := jsoniter.Marshal(user)

		mockUsers.EXPECT().
			RegisterUser(gomock.Any()).
			Return(http.StatusBadRequest, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))

		logger := zerolog.Nop()
		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.PostApiRegisterUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("Login/success", func(t *testing.T) {
		user := models.User{Username: "test"}
		body, _ := jsoniter.Marshal(user)

		mockUsers.EXPECT().
			LoginUser(gomock.Any()).
			Return(&models.TokenResponse{AccessToken: "token"}, http.StatusOK, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiLoginUser(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Login/bad body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		w := httptest.NewRecorder()

		tr.PostApiLoginUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Login/invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer([]byte("bad json")))
		w := httptest.NewRecorder()

		tr.PostApiLoginUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Login/app error", func(t *testing.T) {
		user := models.User{Username: "test"}
		body, _ := jsoniter.Marshal(user)

		mockUsers.EXPECT().
			LoginUser(gomock.Any()).
			Return(nil, http.StatusUnauthorized, errors.New("unauthorized"))

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiLoginUser(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})
}
