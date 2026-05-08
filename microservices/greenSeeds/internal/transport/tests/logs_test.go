package transport_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/mocks"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/transport"
)

func TestLogsTransport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogs := mocks.NewMockILogsApp(ctrl)

	tr := &transport.Transport{
		Logs: mockLogs,
	}

	// ---------- SUCCESS ----------
	t.Run("GetLogs/success", func(t *testing.T) {
		now := time.Now()

		mockLogs.EXPECT().
			GetLogs(gomock.Any()).
			Return([]models.Log{
				{
					Id:  1,
					Dt:  now,
					Lvl: "info",
					Msg: "test log",
				},
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/logs/get?search=test&level=info&limit=10&offset=0&date_from=2024-01-01&date_to=2024-12-31", nil)
		w := httptest.NewRecorder()

		tr.GetApiLogsGet(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	// ---------- EMPTY RESULT ----------
	t.Run("GetLogs/not found", func(t *testing.T) {
		mockLogs.EXPECT().
			GetLogs(gomock.Any()).
			Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/logs/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiLogsGet(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})

	// ---------- APP ERROR ----------
	t.Run("GetLogs/error", func(t *testing.T) {
		mockLogs.EXPECT().
			GetLogs(gomock.Any()).
			Return(nil, assertError())

		req := httptest.NewRequest(http.MethodGet, "/api/logs/get", nil)
		w := httptest.NewRecorder()

		tr.GetApiLogsGet(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})

	// ---------- BAD DATE ----------
	t.Run("GetLogs/bad date", func(t *testing.T) {
		mockLogs.EXPECT().
			GetLogs(gomock.Any()).
			Return([]models.Log{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/logs/get?date_from=bad-date&date_to=bad-date", nil)
		w := httptest.NewRecorder()

		tr.GetApiLogsGet(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})
}

// маленький helper чтобы не тянуть errors.New в каждом тесте
func assertError() error {
	return &testError{}
}

type testError struct{}

func (e *testError) Error() string {
	return "test error"
}
