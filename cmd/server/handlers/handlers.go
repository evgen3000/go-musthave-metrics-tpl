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
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Storage storage.MetricsStorage
}

func NewHandler(storage storage.MetricsStorage) *Handler {
	return &Handler{Storage: storage}
}

func (h *Handler) Ping(rw http.ResponseWriter, _ *http.Request) {
	if h.Storage.StorageType() == "db" {
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) HomeHandler(rw http.ResponseWriter, _ *http.Request) {
	var body strings.Builder
	body.WriteString("<h4>Gauges</h4>")
	for gaugeName, value := range h.Storage.GetAllGauges(context.Background()) {
		body.WriteString(gaugeName + ": " + strconv.FormatFloat(value, 'f', -1, 64) + "</br>")
	}
	body.WriteString("<h4>Counters</h4>")

	for counterName, value := range h.Storage.GetAllCounters(context.Background()) {
		body.WriteString(counterName + ": " + strconv.FormatInt(value, 10) + "</br>")
	}
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, err := rw.Write([]byte(body.String()))
	if err != nil {
		http.Error(rw, "Write failed: %v", http.StatusBadRequest)
	}
}

func (h *Handler) UpdateMetricHandlerJSON(rw http.ResponseWriter, r *http.Request) {
	var body dto.MetricsDTO
	rw.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	switch body.MType {
	case dto.MetricTypeCounter:
		h.Storage.IncrementCounter(context.Background(), body.ID, *body.Delta)
		value, _ := h.Storage.GetCounter(context.Background(), body.ID)

		jsonBody, err := json.Marshal(dto.MetricsDTO{ID: body.ID, MType: body.MType, Delta: &value})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}
		_, err = rw.Write(jsonBody)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}
		rw.WriteHeader(http.StatusOK)
		return
	case dto.MetricTypeGauge:
		h.Storage.SetGauge(context.Background(), body.ID, *body.Value)
		value, _ := h.Storage.GetGauge(context.Background(), body.ID)

		jsonBody, err := json.Marshal(dto.MetricsDTO{ID: body.ID, MType: body.MType, Value: &value})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}
		_, err = rw.Write(jsonBody)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}
		return
	default:
		http.Error(rw, "Bad request", http.StatusBadRequest)
		return
	}
}

func (h *Handler) UpdateMetricHandlerText(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	switch metricType {
	case dto.MetricTypeCounter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(rw, "Bad request", http.StatusBadRequest)
		}
		h.Storage.IncrementCounter(context.Background(), metricName, value)
	case dto.MetricTypeGauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(rw, "Bad request", http.StatusBadRequest)
		}
		h.Storage.SetGauge(context.Background(), metricName, value)
	default:
		http.Error(rw, "Bad request", http.StatusBadRequest)
	}

	rw.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMetricHandlerJSON(rw http.ResponseWriter, r *http.Request) {
	var body dto.MetricsDTO
	rw.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if body.MType != dto.MetricTypeGauge && body.MType != dto.MetricTypeCounter {
		http.Error(rw, "Invalid metric type", http.StatusBadRequest)
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
			http.Error(rw, "Write failed", http.StatusBadRequest)
			return
		}
	} else {
		value, exists := h.Storage.GetCounter(context.Background(), body.ID)
		if !exists {
			http.Error(rw, "Metric not found", http.StatusNotFound)
			return
		}
		jsonBody, err := json.Marshal(dto.MetricsDTO{ID: body.ID, MType: body.MType, Delta: &value})
		if err != nil {
			http.Error(rw, "Json write failed:", http.StatusBadRequest)
			return
		}

		_, err = rw.Write(jsonBody)
		if err != nil {
			http.Error(rw, "Write failed", http.StatusBadRequest)
			return
		}
	}
}

func (h *Handler) GetMetricHandlerText(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if metricType != dto.MetricTypeGauge && metricType != dto.MetricTypeCounter {
		http.Error(rw, "Invalid metric type", http.StatusBadRequest)
	}

	if metricType == dto.MetricTypeGauge {
		value, exists := h.Storage.GetGauge(context.Background(), metricName)
		if !exists {
			http.Error(rw, "Metric not found", http.StatusNotFound)
		}
		rw.Header().Set("Content-Type", "text/plain")
		_, err := rw.Write([]byte(strconv.FormatFloat(value, 'f', -1, 64)))
		if err != nil {
			http.Error(rw, "Write failed", http.StatusBadRequest)
		}
	} else {
		value, exists := h.Storage.GetCounter(context.Background(), metricName)
		if !exists {
			http.Error(rw, "Metric not found", http.StatusNotFound)
		}
		rw.Header().Set("Content-Type", "text/plain")
		_, err := rw.Write([]byte(strconv.FormatInt(value, 10)))
		if err != nil {
			http.Error(rw, "Write failed", http.StatusBadRequest)
		}
	}
}

func (h *Handler) UpdateMetrics(rw http.ResponseWriter, r *http.Request) {
	var body []dto.MetricsDTO
	rw.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	h.Storage.SetMetrics(ctx, body)
	defer cancel()
	rw.WriteHeader(http.StatusOK)
}
