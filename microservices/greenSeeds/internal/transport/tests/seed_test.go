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

func TestSeedsTransport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSeeds := mocks.NewMockISeedsApp(ctrl)

	tr := &transport.Transport{
		Seeds: mockSeeds,
	}

	// ---------- ADD ----------
	t.Run("Add/success", func(t *testing.T) {
		s := models.Seeds{Seed: "corn"}
		body, _ := jsoniter.Marshal(s)

		mockSeeds.EXPECT().
			AddSeed(gomock.Any()).
			Return(s, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/seeds/add", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiSeedAdd(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Add/invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/seeds/add", bytes.NewBuffer([]byte("bad json")))
		w := httptest.NewRecorder()

		tr.PostApiSeedAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Add/app error", func(t *testing.T) {
		s := models.Seeds{Seed: "corn"}
		body, _ := jsoniter.Marshal(s)

		mockSeeds.EXPECT().
			AddSeed(gomock.Any()).
			Return(models.Seeds{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/api/seeds/add", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiSeedAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- GET ----------
	t.Run("Get/success", func(t *testing.T) {
		mockSeeds.EXPECT().
			GetSeeds().
			Return([]models.Seeds{{Seed: "corn"}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/seeds/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiSeedGet(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Get/not found", func(t *testing.T) {
		mockSeeds.EXPECT().
			GetSeeds().
			Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/seeds/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiSeedGet(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	// ---------- GET BY ID ----------
	t.Run("GetById/success", func(t *testing.T) {
		mockSeeds.EXPECT().
			GetSeedBySeed("corn").
			Return(models.Seeds{Seed: "corn"}, nil)

		r := chi.NewRouter()
		r.Get("/{seed}", tr.GetApiSeedGetSeed)

		req := httptest.NewRequest(http.MethodGet, "/corn", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	// ---------- GET WITH BUNKERS ----------
	t.Run("GetWithBunkers/success", func(t *testing.T) {
		mockSeeds.EXPECT().
			GetSeedWithBunkers("corn").
			Return([]models.SeedsWithBunker{
				{
					Seed:         "corn",
					SeedRu:       "кукуруза",
					MinDensity:   10,
					MaxDensity:   20,
					TankCapacity: 100,
					Bunker:       1,
					Amount:       50,
				},
			}, nil)

		r := chi.NewRouter()
		r.Get("/{seed}", tr.GetApiSeedWithBunkers)

		req := httptest.NewRequest(http.MethodGet, "/corn", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	// ---------- UPDATE ----------
	t.Run("Update/success", func(t *testing.T) {
		s := models.Seeds{Seed: "corn"}
		body, _ := jsoniter.Marshal(s)

		mockSeeds.EXPECT().
			UpdateSeed(gomock.Any()).
			Return(s, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/seeds/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiSeedUpdate(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Update/error", func(t *testing.T) {
		s := models.Seeds{Seed: "corn"}
		body, _ := jsoniter.Marshal(s)

		mockSeeds.EXPECT().
			UpdateSeed(gomock.Any()).
			Return(models.Seeds{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPut, "/api/seeds/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiSeedUpdate(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- DELETE ----------
	t.Run("Delete/success", func(t *testing.T) {
		mockSeeds.EXPECT().
			DeleteSeed("corn").
			Return(true, nil)

		r := chi.NewRouter()

		// добавляем логгер в контекст
		logger := zerolog.Nop()
		r.With(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				ctx = context.WithValue(ctx, log.CtxKey, logger)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}).Delete("/{seed}", tr.DeleteApiSeedDelete)

		req := httptest.NewRequest(http.MethodDelete, "/corn", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", w.Code)
		}
	})

	t.Run("Delete/fail", func(t *testing.T) {
		mockSeeds.EXPECT().
			DeleteSeed("corn").
			Return(false, nil)

		r := chi.NewRouter()
		logger := zerolog.Nop()

		r.With(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), log.CtxKey, logger)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}).Delete("/{seed}", tr.DeleteApiSeedDelete)

		req := httptest.NewRequest(http.MethodDelete, "/corn", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}
