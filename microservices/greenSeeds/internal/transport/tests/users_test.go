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

	"github.com/Impisigmatus/service_core/log"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/mocks"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/transport"
	"github.com/rs/zerolog"
)

func TestUsersTransport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsers := mocks.NewMockIUsersApp(ctrl)

	tr := &transport.Transport{
		Users: mockUsers,
	}

	userID := int64(1)
	username := "john"

	// ---------- GET BY ID ----------
	t.Run("GetById/success", func(t *testing.T) {
		u := &models.User{
			Id:       &userID,
			Username: username,
		}

		mockUsers.EXPECT().
			GetUserById("1").
			Return(u, nil)

		r := chi.NewRouter()
		r.Get("/{user_id}", tr.GetApiUserGetUsername)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("GetById/error", func(t *testing.T) {
		mockUsers.EXPECT().
			GetUserById("1").
			Return(nil, errors.New("fail"))

		r := chi.NewRouter()
		r.Get("/{user_id}", tr.GetApiUserGetUsername)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- GET ALL ----------
	t.Run("GetAll/success", func(t *testing.T) {
		mockUsers.EXPECT().
			CheckAllUsers().
			Return([]models.User{{Username: username}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/users/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiCheckAllUsers(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("GetAll/error", func(t *testing.T) {
		mockUsers.EXPECT().
			CheckAllUsers().
			Return(nil, errors.New("fail"))

		req := httptest.NewRequest(http.MethodGet, "/api/users/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiCheckAllUsers(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- CHANGE PASSWORD ----------
	t.Run("ChangePassword/success", func(t *testing.T) {
		body := `{"id":1,"old_password":"old","new_password":"new"}`

		mockUsers.EXPECT().
			ChangePassword(gomock.Any()).
			Return(true, nil)

		logger := zerolog.Nop()

		req := httptest.NewRequest(http.MethodPut, "/api/users/change-password", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		tr.PutApiChangePassword(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}
	})

	t.Run("ChangePassword/not found", func(t *testing.T) {
		body := `{"id":1,"old_password":"old","new_password":"new"}`

		mockUsers.EXPECT().
			ChangePassword(gomock.Any()).
			Return(false, nil)

		logger := zerolog.Nop()

		req := httptest.NewRequest(http.MethodPut, "/api/users/change-password", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		tr.PutApiChangePassword(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("ChangePassword/error", func(t *testing.T) {
		body := `{"id":1,"old_password":"old","new_password":"new"}`

		mockUsers.EXPECT().
			ChangePassword(gomock.Any()).
			Return(false, errors.New("fail"))

		logger := zerolog.Nop()

		req := httptest.NewRequest(http.MethodPut, "/api/users/change-password", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), log.CtxKey, logger)
		req = req.WithContext(ctx)

		tr.PutApiChangePassword(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- DELETE ----------
	t.Run("Delete/success", func(t *testing.T) {
		mockUsers.EXPECT().
			RemoveUser("john").
			Return(true, nil)

		r := chi.NewRouter()
		logger := zerolog.Nop()

		r.With(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), log.CtxKey, logger)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}).Delete("/{username}", tr.DeleteApiRemoveUser)

		req := httptest.NewRequest(http.MethodDelete, "/john", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}
	})

	t.Run("Delete/not found", func(t *testing.T) {
		mockUsers.EXPECT().
			RemoveUser("john").
			Return(false, nil)

		r := chi.NewRouter()
		logger := zerolog.Nop()

		r.With(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), log.CtxKey, logger)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}).Delete("/{username}", tr.DeleteApiRemoveUser)

		req := httptest.NewRequest(http.MethodDelete, "/john", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	// ---------- UPDATE ----------
	t.Run("Update/success", func(t *testing.T) {
		u := models.User{Username: username}
		body, _ := jsoniter.Marshal(u)

		mockUsers.EXPECT().
			Update(gomock.Any()).
			Return(true, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/users/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiUpdateUser(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}
	})

	t.Run("Update/not found", func(t *testing.T) {
		u := models.User{Username: username}
		body, _ := jsoniter.Marshal(u)

		mockUsers.EXPECT().
			Update(gomock.Any()).
			Return(false, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/users/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiUpdateUser(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	t.Run("Update/error", func(t *testing.T) {
		u := models.User{Username: username}
		body, _ := jsoniter.Marshal(u)

		mockUsers.EXPECT().
			Update(gomock.Any()).
			Return(false, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPut, "/api/users/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiUpdateUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}
