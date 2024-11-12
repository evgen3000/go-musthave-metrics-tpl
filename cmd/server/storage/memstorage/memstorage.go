package memstorage

import (
	"context"
	"log"

	"evgen3000/go-musthave-metrics-tpl.git/internal/dto"
)

type MemStorage struct {
	Gauges   map[string]float64 `json:"gauges"`
	Counters map[string]int64   `json:"counters"`
}

func (m *MemStorage) StorageType() string {
	return "ms"
}

func (m *MemStorage) SetMetrics(ctx context.Context, metrics []dto.MetricsDTO) {
	for _, metric := range metrics {
		if metric.MType == dto.MetricTypeCounter {
			m.IncrementCounter(ctx, metric.ID, *metric.Delta)
		} else if metric.MType == dto.MetricTypeGauge {
			m.SetGauge(ctx, metric.ID, *metric.Value)
		} else {
			log.Printf("Unknown metric type: %s", metric.MType)
		}
	}
}

func (m *MemStorage) SetGauge(_ context.Context, metricName string, value float64) {
	m.Gauges[metricName] = value
}

func (m *MemStorage) IncrementCounter(_ context.Context, metricName string, value int64) {
	m.Counters[metricName] += value
}

func (m *MemStorage) GetGauge(_ context.Context, metricName string) (float64, bool) {
	value, exists := m.Gauges[metricName]
	return value, exists
}

func (m *MemStorage) GetCounter(_ context.Context, metricName string) (int64, bool) {
	value, exists := m.Counters[metricName]
	return value, exists
}

func (m *MemStorage) GetAllGauges(_ context.Context) map[string]float64 {
	return m.Gauges
}

func (m *MemStorage) GetAllCounters(_ context.Context) map[string]int64 {
	return m.Counters
}
