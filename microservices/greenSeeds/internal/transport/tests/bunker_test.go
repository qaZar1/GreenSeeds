package transport_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/mocks"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/transport"
)

func TestBunkerTransport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBunkers := mocks.NewMockIBunkersApp(ctrl)

	tr := &transport.Transport{
		Bunkers: mockBunkers,
	}

	// ---------- ADD ----------
	t.Run("Add/success", func(t *testing.T) {
		b := models.Bunkers{Bunker: 1, Distance: 100}
		body, _ := jsoniter.Marshal(b)

		mockBunkers.EXPECT().
			AddBunker(gomock.Any()).
			Return(models.Bunkers{Bunker: 1, Distance: 100}, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/bunkers/add", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiBunkerAdd(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Add/invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/bunkers/add", bytes.NewBuffer([]byte("bad json")))
		w := httptest.NewRecorder()

		tr.PostApiBunkerAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Add/app error", func(t *testing.T) {
		b := models.Bunkers{Bunker: 1, Distance: 100}
		body, _ := jsoniter.Marshal(b)

		mockBunkers.EXPECT().
			AddBunker(gomock.Any()).
			Return(models.Bunkers{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/api/bunkers/add", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiBunkerAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- GET ALL ----------
	t.Run("Get/success", func(t *testing.T) {
		mockBunkers.EXPECT().
			GetBunkers().
			Return([]models.Bunkers{{Bunker: 1, Distance: 100}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/bunkers/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiBunkerGet(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Get/not found", func(t *testing.T) {
		mockBunkers.EXPECT().
			GetBunkers().
			Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/bunkers/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiBunkerGet(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	// ---------- GET FOR PLACEMENT ----------
	t.Run("GetForPlacement/success", func(t *testing.T) {
		mockBunkers.EXPECT().
			GetBunkersForPlacement().
			Return([]models.Bunkers{{Bunker: 1, Distance: 100}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/bunkers/getForPlacement", nil)
		w := httptest.NewRecorder()

		tr.GetApiBunkerGetForPlacement(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	// ---------- GET BY ID ----------
	t.Run("GetById/success", func(t *testing.T) {
		mockBunkers.EXPECT().
			GetBunkersById("1").
			Return(models.Bunkers{Bunker: 1, Distance: 100}, nil)

		r := chi.NewRouter()
		r.Get("/{bunker}", tr.GetApiBunkerGetId)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("GetById/error", func(t *testing.T) {
		mockBunkers.EXPECT().
			GetBunkersById("1").
			Return(models.Bunkers{}, errors.New("fail"))

		r := chi.NewRouter()
		r.Get("/{bunker}", tr.GetApiBunkerGetId)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- UPDATE ----------
	t.Run("Update/success", func(t *testing.T) {
		b := models.Bunkers{Bunker: 1, Distance: 200}
		body, _ := jsoniter.Marshal(b)

		mockBunkers.EXPECT().
			UpdateBunker(gomock.Any()).
			Return(models.Bunkers{Bunker: 1, Distance: 200}, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/bunkers/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiBunkerUpdate(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Update/error", func(t *testing.T) {
		b := models.Bunkers{Bunker: 1, Distance: 200}
		body, _ := jsoniter.Marshal(b)

		mockBunkers.EXPECT().
			UpdateBunker(gomock.Any()).
			Return(models.Bunkers{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPut, "/api/bunkers/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiBunkerUpdate(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- DELETE ----------
	t.Run("Delete/success", func(t *testing.T) {
		mockBunkers.EXPECT().
			DeleteBunker("1").
			Return(true, nil)

		r := chi.NewRouter()
		r.Delete("/{bunker}", tr.DeleteApiBunkerDelete)

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}
	})

	t.Run("Delete/fail", func(t *testing.T) {
		mockBunkers.EXPECT().
			DeleteBunker("1").
			Return(false, nil)

		r := chi.NewRouter()
		r.Delete("/{bunker}", tr.DeleteApiBunkerDelete)

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}
