package transport_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Impisigmatus/service_core/log"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/mocks"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/transport"
)

func TestReceiptsTransport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceipts := mocks.NewMockIReceiptsApp(ctrl)

	tr := &transport.Transport{
		Receipts: mockReceipts,
	}

	now := time.Now()
	id := int64(1)

	// ---------- ADD ----------
	t.Run("Add/success", func(t *testing.T) {
		r := models.Receipts{
			Receipt:     &id,
			Seed:        "seed",
			SeedRu:      "семя",
			Gcode:       "G1",
			Updated:     &now,
			Description: "desc",
		}

		body, _ := jsoniter.Marshal(r)

		mockReceipts.EXPECT().
			AddReceipts(gomock.Any()).
			Return(r, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/receipts/add", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiReceiptsAdd(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Add/invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/receipts/add", bytes.NewBuffer([]byte("bad json")))
		w := httptest.NewRecorder()

		tr.PostApiReceiptsAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Add/app error", func(t *testing.T) {
		r := models.Receipts{Seed: "seed"}
		body, _ := jsoniter.Marshal(r)

		mockReceipts.EXPECT().
			AddReceipts(gomock.Any()).
			Return(models.Receipts{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/api/receipts/add", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiReceiptsAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- GET ----------
	t.Run("Get/success", func(t *testing.T) {
		mockReceipts.EXPECT().
			GetReceipts().
			Return([]models.Receipts{{Seed: "seed"}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/receipts/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiReceiptsGet(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Get/not found", func(t *testing.T) {
		mockReceipts.EXPECT().
			GetReceipts().
			Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/receipts/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiReceiptsGet(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	// ---------- GET BY ID ----------
	t.Run("GetById/success", func(t *testing.T) {
		mockReceipts.EXPECT().
			GetReceiptsByReceipt(1).
			Return(models.Receipts{Seed: "seed"}, nil)

		r := chi.NewRouter()
		r.Get("/{receipt}", tr.GetApiReceiptsGetReceipt)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("GetById/bad param", func(t *testing.T) {
		r := chi.NewRouter()
		r.Get("/{receipt}", tr.GetApiReceiptsGetReceipt)

		req := httptest.NewRequest(http.MethodGet, "/abc", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- UPDATE ----------
	t.Run("Update/success", func(t *testing.T) {
		rp := models.Receipts{Seed: "seed"}
		body, _ := jsoniter.Marshal(rp)

		mockReceipts.EXPECT().
			UpdateReceipts(gomock.Any()).
			Return(rp, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/receipts/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiReceiptsUpdate(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Update/error", func(t *testing.T) {
		rp := models.Receipts{Seed: "seed"}
		body, _ := jsoniter.Marshal(rp)

		mockReceipts.EXPECT().
			UpdateReceipts(gomock.Any()).
			Return(models.Receipts{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPut, "/api/receipts/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiReceiptsUpdate(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- DELETE ----------
	t.Run("Delete/success", func(t *testing.T) {
		mockReceipts.EXPECT().
			DeleteReceipts(1).
			Return(true, nil)

		r := chi.NewRouter()
		r.Delete("/{receipt}", tr.DeleteApiReceiptsDelete)

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)

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
		r.Delete("/{receipt}", tr.DeleteApiReceiptsDelete)

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Delete/fail", func(t *testing.T) {
		mockReceipts.EXPECT().
			DeleteReceipts(1).
			Return(false, nil)

		r := chi.NewRouter()
		r.Delete("/{receipt}", tr.DeleteApiReceiptsDelete)

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)

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
