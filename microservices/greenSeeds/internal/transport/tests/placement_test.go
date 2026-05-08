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

func TestPlacementTransport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPlacements := mocks.NewMockIPlacementsApp(ctrl)

	tr := &transport.Transport{
		Placements: mockPlacements,
	}

	// ---------- ADD ----------
	t.Run("Add/success", func(t *testing.T) {
		p := models.Placement{Bunker: 1}
		body, _ := jsoniter.Marshal(p)

		mockPlacements.EXPECT().
			AddPlacement(gomock.Any()).
			Return(p, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/placement/add", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiPlacementAdd(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Add/invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/placement/add", bytes.NewBuffer([]byte("bad json")))
		w := httptest.NewRecorder()

		tr.PostApiPlacementAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Add/error", func(t *testing.T) {
		p := models.Placement{Bunker: 1}
		body, _ := jsoniter.Marshal(p)

		mockPlacements.EXPECT().
			AddPlacement(gomock.Any()).
			Return(models.Placement{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/api/placement/add", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiPlacementAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- GET ----------
	t.Run("Get/success", func(t *testing.T) {
		mockPlacements.EXPECT().
			GetPlacements().
			Return([]models.Placement{{Bunker: 1}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/placement/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiPlacementGet(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Get/not found", func(t *testing.T) {
		mockPlacements.EXPECT().
			GetPlacements().
			Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/placement/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiPlacementGet(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	// ---------- GET BY BUNKER ----------
	t.Run("GetByBunker/success", func(t *testing.T) {
		mockPlacements.EXPECT().
			GetPlacementByBunker("1").
			Return(models.Placement{Bunker: 1}, nil)

		r := chi.NewRouter()
		r.Get("/{bunker}", tr.GetApiPlacementGetBunker)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("GetByBunker/error", func(t *testing.T) {
		mockPlacements.EXPECT().
			GetPlacementByBunker("1").
			Return(models.Placement{}, errors.New("fail"))

		r := chi.NewRouter()
		r.Get("/{bunker}", tr.GetApiPlacementGetBunker)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- UPDATE ----------
	t.Run("Update/success", func(t *testing.T) {
		p := models.Placement{Bunker: 1}
		body, _ := jsoniter.Marshal(p)

		mockPlacements.EXPECT().
			UpdatePlacement(gomock.Any()).
			Return(p, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/placement/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiPlacementUpdate(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Update/error", func(t *testing.T) {
		p := models.Placement{Bunker: 1}
		body, _ := jsoniter.Marshal(p)

		mockPlacements.EXPECT().
			UpdatePlacement(gomock.Any()).
			Return(models.Placement{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPut, "/api/placement/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiPlacementUpdate(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- DELETE ----------
	t.Run("Delete/success", func(t *testing.T) {
		mockPlacements.EXPECT().
			DeletePlacement("1").
			Return(true, nil)

		r := chi.NewRouter()
		r.Delete("/{bunker}", tr.DeleteApiPlacementDelete)

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}
	})

	t.Run("Delete/fail", func(t *testing.T) {
		mockPlacements.EXPECT().
			DeletePlacement("1").
			Return(false, nil)

		r := chi.NewRouter()
		r.Delete("/{bunker}", tr.DeleteApiPlacementDelete)

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- FILL ----------
	t.Run("Fill/success", func(t *testing.T) {
		fp := models.FillPlacement{
			Seed:    "corn",
			Percent: 50,
		}
		body, _ := jsoniter.Marshal(fp)

		mockPlacements.EXPECT().
			FillPlacement(gomock.Any()).
			Return(models.Placement{Bunker: 1}, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/placement/fill", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiPlacementFill(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Fill/invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/placement/fill", bytes.NewBuffer([]byte("bad json")))
		w := httptest.NewRecorder()

		tr.PutApiPlacementFill(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Fill/app error", func(t *testing.T) {
		fp := models.FillPlacement{
			Seed:    "corn",
			Percent: 50,
		}
		body, _ := jsoniter.Marshal(fp)

		mockPlacements.EXPECT().
			FillPlacement(gomock.Any()).
			Return(models.Placement{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPut, "/api/placement/fill", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiPlacementFill(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}
