package memstorage

import (
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

func (m *MemStorage) SetMetrics(metrics []dto.MetricsDTO) {
	for _, metric := range metrics {
		if metric.MType == dto.MetricTypeCounter {
			m.IncrementCounter(metric.ID, *metric.Delta)
		} else if metric.MType == dto.MetricTypeGauge {
			m.SetGauge(metric.ID, *metric.Value)
		} else {
			log.Printf("Unknown metric type: %s", metric.MType)
		}
	}
}

func (m *MemStorage) SetGauge(metricName string, value float64) {
	m.Gauges[metricName] = value
}

func (m *MemStorage) IncrementCounter(metricName string, value int64) {
	m.Counters[metricName] += value
}

func (m *MemStorage) GetGauge(metricName string) (float64, bool) {
	value, exists := m.Gauges[metricName]
	return value, exists
}

func (m *MemStorage) GetCounter(metricName string) (int64, bool) {
	value, exists := m.Counters[metricName]
	return value, exists
}

func (m *MemStorage) GetAllGauges() map[string]float64 {
	return m.Gauges
}

func (m *MemStorage) GetAllCounters() map[string]int64 {
	return m.Counters
}
