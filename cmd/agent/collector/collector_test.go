package collector_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/agent/collector"
	"evgen3000/go-musthave-metrics-tpl.git/internal/dto"
	"evgen3000/go-musthave-metrics-tpl.git/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPollRuntimeMetrics(t *testing.T) {
	mockCollector := new(mocks.MockCollector)
	mockCollector.On("CollectMetrics").Return([]dto.MetricsDTO{
		{ID: "TestMetric", MType: dto.MetricTypeGauge, Value: float64Ptr(42)},
	})

	agent := &collector.AgentConfig{
		Collector:      mockCollector,
		MetricsChan:    make(chan []dto.MetricsDTO, 1),
		PollCounter:    0,
		Mu:             sync.Mutex{},
		PollInterval:   time.Millisecond * 50,
		ReportInterval: time.Second * 10,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go agent.PollRuntimeMetrics(ctx)

	time.Sleep(100 * time.Millisecond)

	cancel()

	select {
	case metrics := <-agent.MetricsChan:
		assert.Len(t, metrics, 2) // Должно быть 2 метрики: TestMetric и PollCount
		assert.Equal(t, "TestMetric", metrics[0].ID)
		mockCollector.AssertExpectations(t)
	default:
		t.Fatal("Метрики не были отправлены в канал")
	}
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestPollSystemMetrics(t *testing.T) {
	agent := &collector.AgentConfig{
		MetricsChan:  make(chan []dto.MetricsDTO, 1),
		PollInterval: time.Millisecond * 50,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go agent.PollSystemMetrics(ctx)

	time.Sleep(100 * time.Millisecond)

	cancel()

	select {
	case metrics := <-agent.MetricsChan:
		assert.NotEmpty(t, metrics, "Метрики должны быть собраны")

		expectedMetrics := map[string]bool{
			"TotalMemory": true,
			"FreeMemory":  true,
		}

		for i := 1; i <= 16; i++ { // Предполагаем, что максимум 16 CPU
			expectedMetrics[fmt.Sprintf("CPUutilization%d", i)] = true
		}

		for _, m := range metrics {
			assert.Contains(t, expectedMetrics, m.ID, "Неожиданная метрика: %s", m.ID)
		}
	default:
		t.Fatal("Метрики не были отправлены в канал")
	}
}

func TestSendMetrics(t *testing.T) {
	mockHTTPClient := new(mocks.MockHTTPClient)
	mockHTTPClient.On("SendMetrics", mock.Anything).Return(nil)

	agent := &collector.AgentConfig{
		HTTPClient:  mockHTTPClient,
		MetricsChan: make(chan []dto.MetricsDTO, 1),
		RateLimit:   1,
	}

	metrics := []dto.MetricsDTO{
		{ID: "TestMetric", MType: dto.MetricTypeGauge, Value: float64Ptr(42)},
	}
	agent.MetricsChan <- metrics

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go agent.SendMetrics(ctx)

	time.Sleep(100 * time.Millisecond)

	mockHTTPClient.AssertCalled(t, "SendMetrics", mock.Anything)
}
