package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/handlers"
	"evgen3000/go-musthave-metrics-tpl.git/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPing(t *testing.T) {
	mockStorage := new(mocks.MockMetricsStorage)
	mockStorage.On("StorageType").Return("db")

	handler := handlers.NewHandler(mockStorage, nil)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	handler.Ping(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockStorage.AssertExpectations(t)
}

func TestHomeHandler(t *testing.T) {
	mockStorage := new(mocks.MockMetricsStorage)
	mockStorage.On("GetAllGauges", mock.Anything).Return(map[string]float64{
		"gauge1": 42.42,
		"gauge2": 24.24,
	})
	mockStorage.On("GetAllCounters", mock.Anything).Return(map[string]int64{
		"counter1": 100,
		"counter2": 200,
	})

	handler := handlers.NewHandler(mockStorage, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.HomeHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "gauge1: 42.42")
	assert.Contains(t, rec.Body.String(), "counter1: 100")
	mockStorage.AssertExpectations(t)
}

func TestUpdateMetricHandlerJSON(t *testing.T) {
	mockStorage := new(mocks.MockMetricsStorage)
	mockStorage.On("IncrementCounter", mock.Anything, "counter1", int64(10))
	mockStorage.On("GetCounter", mock.Anything, "counter1").Return(int64(10), true)
	mockStorage.On("SetGauge", mock.Anything, "gauge1", 42.42)
	mockStorage.On("GetGauge", mock.Anything, "gauge1").Return(42.42, true)

	handler := handlers.NewHandler(mockStorage, nil)

	// Test counter update
	counterBody := `{"id": "counter1", "type": "counter", "delta": 10}`
	req := httptest.NewRequest(http.MethodPost, "/update/", strings.NewReader(counterBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.UpdateMetricHandlerJSON(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"delta":10`)

	// Test gauge update
	gaugeBody := `{"id": "gauge1", "type": "gauge", "value": 42.42}`
	req = httptest.NewRequest(http.MethodPost, "/update/", strings.NewReader(gaugeBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()

	handler.UpdateMetricHandlerJSON(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"value":42.42`)

	mockStorage.AssertExpectations(t)
}

func TestGetMetricHandlerJSON(t *testing.T) {
	mockStorage := new(mocks.MockMetricsStorage)
	mockStorage.On("GetCounter", mock.Anything, "counter1").Return(int64(10), true)
	mockStorage.On("GetGauge", mock.Anything, "gauge1").Return(42.42, true)

	handler := handlers.NewHandler(mockStorage, nil)

	// Test counter get
	counterBody := `{"id": "counter1", "type": "counter"}`
	req := httptest.NewRequest(http.MethodPost, "/value/", strings.NewReader(counterBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.GetMetricHandlerJSON(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"delta":10`)

	// Test gauge get
	gaugeBody := `{"id": "gauge1", "type": "gauge"}`
	req = httptest.NewRequest(http.MethodPost, "/value/", strings.NewReader(gaugeBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()

	handler.GetMetricHandlerJSON(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"value":42.42`)

	mockStorage.AssertExpectations(t)
}
