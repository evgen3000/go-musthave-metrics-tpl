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

func GenerateJSON(m []dto.MetricsDTO) []byte {
	body, err := json.Marshal(m)
	if err != nil {
		log.Fatal("Conversion have errors:", err.Error())
	}
	return body
}

type AgentConfig struct {
	host           string
	pollInterval   time.Duration
	reportInterval time.Duration
	collector      *metrics.Collector
	httpClient     *httpclient.HTTPClient
	rateLimit      int
	metricsChan    chan []dto.MetricsDTO
}

func NewAgent(host string, pollInterval, reportInterval time.Duration, key string, rateLimit int) *AgentConfig {
	return &AgentConfig{
		host:           host,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		collector:      metrics.NewMetricsCollector(),
		httpClient:     httpclient.NewHTTPClient(host, key),
		rateLimit:      rateLimit,
		metricsChan:    make(chan []dto.MetricsDTO, rateLimit),
	}
}

func (a *AgentConfig) Start(ctx context.Context) {
	var wg sync.WaitGroup

	// Горутина для сбора runtime-метрик
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.pollRuntimeMetrics(ctx)
	}()

	// Горутина для сбора системных метрик
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.pollSystemMetrics(ctx)
	}()

	// Горутина для отправки метрик
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.sendMetrics(ctx)
	}()

	wg.Wait()
	fmt.Println("Agent завершил работу.")
}

func (a *AgentConfig) pollRuntimeMetrics(ctx context.Context) {
	ticker := time.NewTicker(a.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Сбор runtime метрик завершен.")
			return
		case <-ticker.C:
			metricsTicker := a.collector.CollectMetrics()
			a.metricsChan <- metricsTicker
		}
	}
}

func (a *AgentConfig) pollSystemMetrics(ctx context.Context) {
	ticker := time.NewTicker(a.pollInterval)
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

			a.metricsChan <- systemMetrics
		}
	}
}

func (a *AgentConfig) sendMetrics(ctx context.Context) {
	sem := make(chan struct{}, a.rateLimit) // Ограничение запросов

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Отправка метрик завершена.")
			return
		case metricsTicker := <-a.metricsChan:
			sem <- struct{}{} // Захват семафора
			go func(data []dto.MetricsDTO) {
				defer func() { <-sem }() // Освобождение семафора
				jsonData := GenerateJSON(data)
				a.httpClient.SendMetrics(jsonData)
				log.Printf("Отправлены метрики: %s\n", string(jsonData))
			}(metricsTicker)
		}
	}
}
