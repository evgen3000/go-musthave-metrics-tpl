package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage"
	"evgen3000/go-musthave-metrics-tpl.git/internal/dto"
	"evgen3000/go-musthave-metrics-tpl.git/internal/workerpool"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Storage    storage.MetricsStorage
	WorkerPool *workerpool.WorkerPool
}

// NewHandler создает новый экземпляр Handler.
func NewHandler(storage storage.MetricsStorage, workerPool *workerpool.WorkerPool) *Handler {
	return &Handler{
		Storage:    storage,
		WorkerPool: workerPool,
	}
}

// Ping проверяет доступность хранилища.
func (h *Handler) Ping(rw http.ResponseWriter, _ *http.Request) {
	if h.Storage.StorageType() == "db" {
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

// HomeHandler отображает все доступные метрики.
func (h *Handler) HomeHandler(rw http.ResponseWriter, _ *http.Request) {
	var body strings.Builder
	body.WriteString("<h4>Gauges</h4>")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for gaugeName, value := range h.Storage.GetAllGauges(ctx) {
		body.WriteString(gaugeName + ": " + strconv.FormatFloat(value, 'f', -1, 64) + "</br>")
	}
	body.WriteString("<h4>Counters</h4>")
	for counterName, value := range h.Storage.GetAllCounters(ctx) {
		body.WriteString(counterName + ": " + strconv.FormatInt(value, 10) + "</br>")
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := rw.Write([]byte(body.String()))
	if err != nil {
		http.Error(rw, "Write failed", http.StatusInternalServerError)
	}
}

// UpdateMetrics обрабатывает обновление метрик через JSON.
func (h *Handler) UpdateMetrics(rw http.ResponseWriter, r *http.Request) {
	var metrics []dto.MetricsDTO
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавляем задачу в пул воркеров
	h.WorkerPool.Submit(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		h.Storage.SetMetrics(ctx, metrics)
	})

	rw.WriteHeader(http.StatusOK)
}

// UpdateMetricHandlerJSON обрабатывает обновление одной метрики через JSON.
func (h *Handler) UpdateMetricHandlerJSON(rw http.ResponseWriter, r *http.Request) {
	var body dto.MetricsDTO
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	switch body.MType {
	case dto.MetricTypeCounter:
		h.WorkerPool.Submit(func() {
			h.Storage.IncrementCounter(context.Background(), body.ID, *body.Delta)
		})
	case dto.MetricTypeGauge:
		h.WorkerPool.Submit(func() {
			h.Storage.SetGauge(context.Background(), body.ID, *body.Value)
		})
	default:
		http.Error(rw, "Bad request", http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// UpdateMetricHandlerText обрабатывает обновление одной метрики через URL параметры.
func (h *Handler) UpdateMetricHandlerText(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	switch metricType {
	case dto.MetricTypeCounter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(rw, "Bad request", http.StatusBadRequest)
			return
		}
		h.WorkerPool.Submit(func() {
			h.Storage.IncrementCounter(context.Background(), metricName, value)
		})
	case dto.MetricTypeGauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(rw, "Bad request", http.StatusBadRequest)
			return
		}
		h.WorkerPool.Submit(func() {
			h.Storage.SetGauge(context.Background(), metricName, value)
		})
	default:
		http.Error(rw, "Bad request", http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// GetMetricHandlerJSON обрабатывает запрос на получение метрики через JSON.
func (h *Handler) GetMetricHandlerJSON(rw http.ResponseWriter, r *http.Request) {
	var body dto.MetricsDTO
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if body.MType == dto.MetricTypeGauge {
		value, exists := h.Storage.GetGauge(context.Background(), body.ID)
		if !exists {
			http.Error(rw, "Metric not found", http.StatusNotFound)
			return
		}
		jsonBody, _ := json.Marshal(dto.MetricsDTO{ID: body.ID, MType: body.MType, Value: &value})
		_, err := rw.Write(jsonBody)
		if err != nil {
			http.Error(rw, "Write failed", http.StatusInternalServerError)
		}
	} else if body.MType == dto.MetricTypeCounter {
		value, exists := h.Storage.GetCounter(context.Background(), body.ID)
		if !exists {
			http.Error(rw, "Metric not found", http.StatusNotFound)
			return
		}
		jsonBody, _ := json.Marshal(dto.MetricsDTO{ID: body.ID, MType: body.MType, Delta: &value})
		_, err := rw.Write(jsonBody)
		if err != nil {
			http.Error(rw, "Write failed", http.StatusInternalServerError)
		}
	} else {
		http.Error(rw, "Invalid metric type", http.StatusBadRequest)
	}
}

// GetMetricHandlerText обрабатывает запрос на получение метрики через URL параметры.
func (h *Handler) GetMetricHandlerText(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if metricType == dto.MetricTypeGauge {
		value, exists := h.Storage.GetGauge(context.Background(), metricName)
		if !exists {
			http.Error(rw, "Metric not found", http.StatusNotFound)
			return
		}
		_, err := rw.Write([]byte(strconv.FormatFloat(value, 'f', -1, 64)))
		if err != nil {
			http.Error(rw, "Write failed", http.StatusInternalServerError)
		}
	} else if metricType == dto.MetricTypeCounter {
		value, exists := h.Storage.GetCounter(context.Background(), metricName)
		if !exists {
			http.Error(rw, "Metric not found", http.StatusNotFound)
			return
		}
		_, err := rw.Write([]byte(strconv.FormatInt(value, 10)))
		if err != nil {
			http.Error(rw, "Write failed", http.StatusInternalServerError)
		}
	} else {
		http.Error(rw, "Invalid metric type", http.StatusBadRequest)
	}
}
