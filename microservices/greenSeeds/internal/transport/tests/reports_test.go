package transport_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/mocks"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/transport"
)

func TestReportsTransport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReports := mocks.NewMockIReportsApp(ctrl)

	tr := &transport.Transport{
		Reports: mockReports,
	}

	now := time.Now()
	id := int64(1)

	// ---------- ADD ----------
	t.Run("Add/success", func(t *testing.T) {
		rep := models.Reports{
			Id:      &id,
			Shift:   1,
			Number:  1,
			Recipe:  1,
			Turn:    1,
			Dt:      &now,
			Success: true,
		}

		body, _ := jsoniter.Marshal(rep)

		mockReports.EXPECT().
			AddReport(gomock.Any()).
			Return(rep, nil)

		req := httptest.NewRequest(http.MethodPost, "/api/reports/add", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiReportsAdd(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Add/invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/reports/add", bytes.NewBuffer([]byte("bad json")))
		w := httptest.NewRecorder()

		tr.PostApiReportsAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Add/app error", func(t *testing.T) {
		rep := models.Reports{Shift: 1}
		body, _ := jsoniter.Marshal(rep)

		mockReports.EXPECT().
			AddReport(gomock.Any()).
			Return(models.Reports{}, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/api/reports/add", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PostApiReportsAdd(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- GET ----------
	t.Run("Get/success", func(t *testing.T) {
		mockReports.EXPECT().
			GetReports().
			Return([]models.Reports{{Shift: 1}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/reports/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiReports(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Get/not found", func(t *testing.T) {
		mockReports.EXPECT().
			GetReports().
			Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/reports/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiReports(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	// ---------- GET BY ID ----------
	t.Run("GetById/success", func(t *testing.T) {
		mockReports.EXPECT().
			GetReportsByReport("1").
			Return(models.Reports{Shift: 1}, nil)

		r := chi.NewRouter()
		r.Get("/{id}", tr.GetApiReportsById)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("GetById/error", func(t *testing.T) {
		mockReports.EXPECT().
			GetReportsByReport("1").
			Return(models.Reports{}, errors.New("fail"))

		r := chi.NewRouter()
		r.Get("/{id}", tr.GetApiReportsById)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- UPDATE ----------
	t.Run("Update/success", func(t *testing.T) {
		rep := models.Reports{Shift: 1}
		body, _ := jsoniter.Marshal(rep)

		mockReports.EXPECT().
			UpdateReport(gomock.Any()).
			Return(true, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/reports/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiReportsUpdate(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("Update/error", func(t *testing.T) {
		rep := models.Reports{Shift: 1}
		body, _ := jsoniter.Marshal(rep)

		mockReports.EXPECT().
			UpdateReport(gomock.Any()).
			Return(false, errors.New("fail"))

		req := httptest.NewRequest(http.MethodPut, "/api/reports/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiReportsUpdate(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	t.Run("Update/not ok", func(t *testing.T) {
		rep := models.Reports{Shift: 1}
		body, _ := jsoniter.Marshal(rep)

		mockReports.EXPECT().
			UpdateReport(gomock.Any()).
			Return(false, nil)

		req := httptest.NewRequest(http.MethodPut, "/api/reports/update", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		tr.PutApiReportsUpdate(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}
