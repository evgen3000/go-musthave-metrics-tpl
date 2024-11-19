package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/router"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/memstorage"
	"github.com/stretchr/testify/assert"
)

//func TestRouter_UpdateMetricHandlerJSON(t *testing.T) {
//	storage := &memstorage.MemStorage{
//		Gauges:   make(map[string]float64),
//		Counters: make(map[string]int64),
//	}
//	chiRouter := router.SetupRouter(storage, "123")
//
//	reqBody := map[string]interface{}{
//		"id":    "testGauge",
//		"type":  "gauge",
//		"value": 42.42,
//	}
//	jsonBody, _ := json.Marshal(reqBody)
//
//	req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(jsonBody))
//	req.Header.Set("Content-Type", "application/json")
//	resp := httptest.NewRecorder()
//	chiRouter.ServeHTTP(resp, req)
//
//	assert.Equal(t, http.StatusOK, resp.Code)
//}

//func TestRouter_UpdateMetricHandlerText(t *testing.T) {
//	storage := &memstorage.MemStorage{
//		Gauges:   make(map[string]float64),
//		Counters: make(map[string]int64),
//	}
//	chiRouter := router.SetupRouter(storage, "123")
//
//	req := httptest.NewRequest(http.MethodPost, "/update/gauge/testGauge/42.42", nil)
//	resp := httptest.NewRecorder()
//	chiRouter.ServeHTTP(resp, req)
//
//	assert.Equal(t, http.StatusOK, resp.Code)
//}

func TestRouter_GetMetricHandlerJSON(t *testing.T) {
	storage := &memstorage.MemStorage{
		Gauges: map[string]float64{
			"testGauge": 42.42,
		},
		Counters: make(map[string]int64),
	}
	chiRouter := router.SetupRouter(storage, "123")

	reqBody := map[string]interface{}{
		"id":   "testGauge",
		"type": "gauge",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	chiRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestRouter_HomeHandler(t *testing.T) {
	storage := &memstorage.MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
	chiRouter := router.SetupRouter(storage, "123")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	chiRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Gauges")
}
