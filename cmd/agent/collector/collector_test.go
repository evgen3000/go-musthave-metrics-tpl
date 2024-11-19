package collector_test

import (
	"context"
	"testing"
	"time"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/agent/collector"
	"evgen3000/go-musthave-metrics-tpl.git/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJSON(t *testing.T) {
	metrics := []dto.MetricsDTO{
		{ID: "testMetric", MType: "gauge", Value: func(v float64) *float64 { return &v }(42.42)},
	}
	jsonData := collector.GenerateJSON(metrics)

	expected := `[{"id":"testMetric","type":"gauge","value":42.42}]`
	assert.JSONEq(t, expected, string(jsonData))
}

func TestAgentConfig_PoolCount(t *testing.T) {
	host := "http://localhost"
	pollInterval := 100 * time.Millisecond
	reportInterval := 200 * time.Millisecond
	agent := collector.NewAgent(host, pollInterval, reportInterval, "123")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	agent.Start(ctx)

	assert.True(t, agent.PoolCount > 0, "Expected PoolCount to be greater than 0")
}
