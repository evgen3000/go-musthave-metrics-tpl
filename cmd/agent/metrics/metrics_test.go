package metrics_test

import (
	"testing"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/agent/metrics"
	"evgen3000/go-musthave-metrics-tpl.git/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestCollectMetrics(t *testing.T) {
	collector := metrics.NewMetricsCollector()

	collectedMetrics := collector.CollectMetrics()

	assert.NotEmpty(t, collectedMetrics, "Collected metrics should not be empty")

	for _, metric := range collectedMetrics {
		assert.NotEmpty(t, metric.ID, "Metric ID should not be empty")
		assert.Equal(t, dto.MetricTypeGauge, metric.MType, "Metric type should be gauge")
		assert.NotNil(t, metric.Value, "Metric value should not be nil")
		assert.Nil(t, metric.Delta, "Metric delta should be nil for gauge metrics")
	}

	// Проверяем наличие ожидаемых метрик
	expectedMetrics := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
		"Sys", "TotalAlloc", "RandomValue",
	}

	// Собираем ID метрик
	metricIDs := make(map[string]bool)
	for _, metric := range collectedMetrics {
		metricIDs[metric.ID] = true
	}

	for _, expectedMetric := range expectedMetrics {
		assert.Contains(t, metricIDs, expectedMetric, "Metric ID %s should be present", expectedMetric)
	}

	for _, metric := range collectedMetrics {
		if metric.ID == "RandomValue" {
			assert.GreaterOrEqual(t, *metric.Value, 0.0, "RandomValue should be >= 0")
			assert.LessOrEqual(t, *metric.Value, 100.0, "RandomValue should be <= 100")
		}
	}
}
