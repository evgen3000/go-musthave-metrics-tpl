package mocks

import (
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
