package mocks

import (
	"context"

	"evgen3000/go-musthave-metrics-tpl.git/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockCollector struct {
	mock.Mock
}

func (m *MockCollector) CollectMetrics() []dto.MetricsDTO {
	args := m.Called()
	return args.Get(0).([]dto.MetricsDTO)
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) SendMetrics(data []byte) {
	m.Called(data)
}

type MockMetricsStorage struct {
	mock.Mock
}

func (m *MockMetricsStorage) StorageType() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockMetricsStorage) GetAllGauges(ctx context.Context) map[string]float64 {
	args := m.Called(ctx)
	return args.Get(0).(map[string]float64)
}

func (m *MockMetricsStorage) GetAllCounters(ctx context.Context) map[string]int64 {
	args := m.Called(ctx)
	return args.Get(0).(map[string]int64)
}

func (m *MockMetricsStorage) IncrementCounter(ctx context.Context, name string, value int64) {
	m.Called(ctx, name, value)
}

func (m *MockMetricsStorage) SetGauge(ctx context.Context, name string, value float64) {
	m.Called(ctx, name, value)
}

func (m *MockMetricsStorage) GetCounter(ctx context.Context, name string) (int64, bool) {
	args := m.Called(ctx, name)
	return args.Get(0).(int64), args.Bool(1)
}

func (m *MockMetricsStorage) GetGauge(ctx context.Context, name string) (float64, bool) {
	args := m.Called(ctx, name)
	return args.Get(0).(float64), args.Bool(1)
}

func (m *MockMetricsStorage) SetMetrics(ctx context.Context, metrics []dto.MetricsDTO) {
	m.Called(ctx, metrics)
}
