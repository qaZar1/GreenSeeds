// go test ./internal/transport/... -coverpkg=./internal/transport
package transport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Impisigmatus/service_core/log"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/mocks"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/transport"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestPostApiAssignmentsAdd(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockIAssignmentsApp(ctrl)

		input := models.Assignments{Shift: 1, Number: 1, Receipt: 10, Amount: 5}

		mockApp.EXPECT().AddAssignment(input).Return(input, nil)

		tr := transport.Transport{Assignments: mockApp}

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		w := httptest.NewRecorder()

		tr.PostApiAssignmentsAdd(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("BadJSON", func(t *testing.T) {
		tr := transport.Transport{}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{bad"))
		w := httptest.NewRecorder()

		tr.PostApiAssignmentsAdd(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("AppError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockIAssignmentsApp(ctrl)
		input := models.Assignments{Shift: 1, Number: 1, Receipt: 10, Amount: 5}

		mockApp.EXPECT().AddAssignment(input).Return(models.Assignments{}, errors.New("err"))

		tr := transport.Transport{Assignments: mockApp}

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		w := httptest.NewRecorder()

		tr.PostApiAssignmentsAdd(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

func TestGetApiAssignmentsGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockIAssignmentsApp(ctrl)

		mockApp.EXPECT().GetAssignments().Return([]models.Assignments{{}}, nil)

		tr := transport.Transport{Assignments: mockApp}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		tr.GetApiAssignmentsGet(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("NilResult", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockIAssignmentsApp(ctrl)

		mockApp.EXPECT().GetAssignments().Return(nil, nil)

		tr := transport.Transport{Assignments: mockApp}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		tr.GetApiAssignmentsGet(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}

func TestGetApiAssignmentsGetAssignment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockIAssignmentsApp(ctrl)

		mockApp.EXPECT().GetAssignmentsByAssignment("1").Return(models.Assignments{}, nil)

		tr := transport.Transport{Assignments: mockApp}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		tr.GetApiAssignmentsGetAssignment(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
}

func TestPutApiAssignmentsUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockIAssignmentsApp(ctrl)

		input := models.Assignments{Shift: 1}

		mockApp.EXPECT().UpdateAssignment(input).Return(input, nil)

		tr := transport.Transport{Assignments: mockApp}

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(body))
		w := httptest.NewRecorder()

		tr.PutApiAssignmentsUpdate(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("BadJSON", func(t *testing.T) {
		tr := transport.Transport{}

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBufferString("{bad"))
		w := httptest.NewRecorder()

		tr.PutApiAssignmentsUpdate(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

func TestDeleteApiAssignmentsDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockIAssignmentsApp(ctrl)

		mockApp.EXPECT().DeleteAssignments("1").Return(true, nil)

		tr := transport.Transport{Assignments: mockApp}

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)

		logger := zerolog.Nop()
		ctx = context.WithValue(ctx, log.CtxKey, logger)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		tr.DeleteApiAssignmentsDelete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
	})

	t.Run("NoLogger", func(t *testing.T) {
		tr := transport.Transport{}

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		w := httptest.NewRecorder()

		tr.DeleteApiAssignmentsDelete(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

func TestGetApiActiveTasks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockIAssignmentsApp(ctrl)

		mockApp.EXPECT().CheckActiveTasks("user").Return([]models.ActiveTask{{}}, nil)

		tr := transport.Transport{Assignments: mockApp}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("user_id", "user")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		tr.GetApiActiveTasks(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
}

func TestGetApiTask(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockIAssignmentsApp(ctrl)

		mockApp.EXPECT().GetTaskById("1").Return(models.Task{Id: 1}, nil)

		tr := transport.Transport{Assignments: mockApp}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		tr.GetApiTask(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := mocks.NewMockIAssignmentsApp(ctrl)

		mockApp.EXPECT().GetTaskById("1").Return(models.Task{}, nil)

		tr := transport.Transport{Assignments: mockApp}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()

		tr.GetApiTask(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}
