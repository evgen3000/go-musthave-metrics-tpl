package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/agent/httpclient"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/agent/metrics"
	"evgen3000/go-musthave-metrics-tpl.git/internal/dto"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type Collector interface {
	CollectMetrics() []dto.MetricsDTO
}

type CollectorImpl struct{}

func GenerateJSON(m []dto.MetricsDTO) []byte {
	body, err := json.Marshal(m)
	if err != nil {
		log.Fatal("Conversion have errors:", err.Error())
	}
	return body
}

type AgentConfig struct {
	host           string
	PollInterval   time.Duration
	ReportInterval time.Duration
	Collector      metrics.Collector
	HTTPClient     httpclient.HTTPClientInterface
	RateLimit      int
	MetricsChan    chan []dto.MetricsDTO
	Mu             sync.Mutex
	PollCounter    int
}

func NewAgent(host string, pollInterval, reportInterval time.Duration, key string, rateLimit int) *AgentConfig {
	return &AgentConfig{
		host:           host,
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		Collector:      metrics.NewMetricsCollector(), // Это работает как интерфейс
		HTTPClient:     httpclient.NewHTTPClient(host, key),
		RateLimit:      rateLimit,
		MetricsChan:    make(chan []dto.MetricsDTO, rateLimit),
		PollCounter:    0,
	}
}

func (a *AgentConfig) Start(ctx context.Context) {
	var wg sync.WaitGroup

	// Горутина для сбора runtime-метрик
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.PollRuntimeMetrics(ctx)
	}()

	// Горутина для сбора системных метрик
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.PollSystemMetrics(ctx)
	}()

	// Горутина для отправки метрик
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.SendMetrics(ctx)
	}()

	wg.Wait()
	fmt.Println("Agent завершил работу.")
}

func (a *AgentConfig) PollRuntimeMetrics(ctx context.Context) {
	pollTicker := time.NewTicker(a.PollInterval)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Сбор runtime метрик завершен.")
			return
		case <-pollTicker.C:
			collectedMetrics := a.Collector.CollectMetrics()

			// Увеличиваем PoolCount
			a.Mu.Lock()
			a.PollCounter++
			poolCount := int64(a.PollCounter)
			a.Mu.Unlock()

			// Добавляем PoolCount к метрикам
			collectedMetrics = append(collectedMetrics, dto.MetricsDTO{ID: "PollCount", MType: "counter", Delta: &poolCount})

			// Отправляем в канал
			a.MetricsChan <- collectedMetrics
			fmt.Printf("Runtime metrics collected: %v\n", collectedMetrics)
		}
	}
}

func (a *AgentConfig) PollSystemMetrics(ctx context.Context) {
	ticker := time.NewTicker(a.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Сбор системных метрик завершен.")
			return
		case <-ticker.C:
			vm, err := mem.VirtualMemory()
			if err != nil {
				fmt.Println("Ошибка получения данных виртуальной памяти:", err)
				continue
			}
			cpuUtilization, err := cpu.Percent(0, true)
			if err != nil {
				fmt.Println("Ошибка получения загрузки CPU:", err)
				continue
			}

			totalMemory := float64(vm.Total)
			freeMemory := float64(vm.Free)

			systemMetrics := []dto.MetricsDTO{
				{ID: "TotalMemory", MType: dto.MetricTypeGauge, Value: &totalMemory},
				{ID: "FreeMemory", MType: dto.MetricTypeGauge, Value: &freeMemory},
			}

			for i, util := range cpuUtilization {
				id := fmt.Sprintf("CPUutilization%d", i+1)
				systemMetrics = append(systemMetrics, dto.MetricsDTO{ID: id, MType: dto.MetricTypeGauge, Value: &util})
			}

			a.MetricsChan <- systemMetrics
		}
	}
}

func (a *AgentConfig) SendMetrics(ctx context.Context) {
	sem := make(chan struct{}, a.RateLimit) // Ограничение запросов

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Отправка метрик завершена.")
			return
		case metricsTicker := <-a.MetricsChan:
			sem <- struct{}{} // Захват семафора
			go func(data []dto.MetricsDTO) {
				defer func() { <-sem }() // Освобождение семафора
				jsonData := GenerateJSON(data)
				a.HTTPClient.SendMetrics(jsonData)
				log.Printf("Отправлены метрики: %s\n", string(jsonData))
			}(metricsTicker)
		}
	}
}
