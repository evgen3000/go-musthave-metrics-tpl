package router_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	r "evgen3000/go-musthave-metrics-tpl.git/cmd/server/router"
	"evgen3000/go-musthave-metrics-tpl.git/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPingRoute(t *testing.T) {
	mockStorage := new(mocks.MockMetricsStorage)
	mockStorage.On("StorageType").Return("db")

	router := r.SetupRouter(mockStorage, "", nil)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockStorage.AssertExpectations(t)
}

func TestHomeHandlerRoute(t *testing.T) {
	mockStorage := new(mocks.MockMetricsStorage)
	mockStorage.On("GetAllGauges", mock.Anything).Return(map[string]float64{
		"gauge1": 42.42,
	})
	mockStorage.On("GetAllCounters", mock.Anything).Return(map[string]int64{
		"counter1": 100,
	})

	router := r.SetupRouter(mockStorage, "", nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "gauge1: 42.42")
	assert.Contains(t, rec.Body.String(), "counter1: 100")
	mockStorage.AssertExpectations(t)
}

func TestUpdateMetricHandlerJSONRoute(t *testing.T) {
	mockStorage := new(mocks.MockMetricsStorage)
	mockStorage.On("IncrementCounter", mock.Anything, "counter1", int64(10))
	mockStorage.On("GetCounter", mock.Anything, "counter1").Return(int64(10), true)

	router := r.SetupRouter(mockStorage, "", nil)

	body := `{"id": "counter1", "type": "counter", "delta": 10}`
	req := httptest.NewRequest(http.MethodPost, "/update/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"delta":10`)
	mockStorage.AssertExpectations(t)
}

func TestGetMetricHandlerJSONRoute(t *testing.T) {
	mockStorage := new(mocks.MockMetricsStorage)
	mockStorage.On("GetCounter", mock.Anything, "counter1").Return(int64(10), true)

	router := r.SetupRouter(mockStorage, "", nil)

	body := `{"id": "counter1", "type": "counter"}`
	req := httptest.NewRequest(http.MethodPost, "/value/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"delta":10`)
	mockStorage.AssertExpectations(t)
}

func TestUpdateMetricsRoute(t *testing.T) {
	mockStorage := new(mocks.MockMetricsStorage)
	mockStorage.On("SetMetrics", mock.Anything, mock.Anything).Return()

	router := r.SetupRouter(mockStorage, "", nil)

	body := `[{"id": "counter1", "type": "counter", "delta": 10}]`
	req := httptest.NewRequest(http.MethodPost, "/updates/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockStorage.AssertExpectations(t)
}
