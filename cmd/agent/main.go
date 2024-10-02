package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type Agent struct {
	pollInterval   time.Duration
	reportInterval time.Duration
	poolCount      int64
}

func NewAgent(poolInterval, reportInterval time.Duration) *Agent {
	return &Agent{
		pollInterval:   poolInterval,
		reportInterval: reportInterval,
		poolCount:      0,
	}
}

func (a *Agent) collectMetrics() map[string]float64 {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)

	metrics := map[string]float64{"Alloc": float64(memStats.Alloc),
		"BuckHashSys":   float64(memStats.BuckHashSys),
		"Frees":         float64(memStats.Frees),
		"GCCPUFraction": memStats.GCCPUFraction,
		"GCSys":         float64(memStats.GCSys),
		"HeapAlloc":     float64(memStats.HeapAlloc),
		"HeapIdle":      float64(memStats.HeapIdle),
		"HeapInuse":     float64(memStats.HeapInuse),
		"HeapObjects":   float64(memStats.HeapObjects),
		"HeapReleased":  float64(memStats.HeapReleased),
		"HeapSys":       float64(memStats.HeapSys),
		"LastGC":        float64(memStats.LastGC),
		"Lookups":       float64(memStats.Lookups),
		"MCacheInuse":   float64(memStats.MCacheInuse),
		"MCacheSys":     float64(memStats.MCacheSys),
		"MSpanInuse":    float64(memStats.MSpanInuse),
		"MSpanSys":      float64(memStats.MSpanSys),
		"Mallocs":       float64(memStats.Mallocs),
		"NextGC":        float64(memStats.NextGC),
		"NumForcedGC":   float64(memStats.NumForcedGC),
		"NumGC":         float64(memStats.NumGC),
		"OtherSys":      float64(memStats.OtherSys),
		"PauseTotalNs":  float64(memStats.PauseTotalNs),
		"StackInuse":    float64(memStats.StackInuse),
		"StackSys":      float64(memStats.StackSys),
		"Sys":           float64(memStats.Sys),
		"TotalAlloc":    float64(memStats.TotalAlloc),
		"RandomValue":   rand.Float64() * 100,
	}

	return metrics
}

func (a *Agent) sendMetrics(metricType, metricName string, value float64) {
	metricValue := strconv.FormatFloat(value, 'f', -1, 64)
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", metricType, metricName, metricValue)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Println("Error of creating request:", err)
		return
	}
	req.Header.Set("Content-type", "text-plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error to sending request:", err)
		return
	}

	defer resp.Body.Close()
	fmt.Printf("Metrics %s (%s) with value %s sent succesfully", metricName, metricType, metricValue)
}

func (a *Agent) start() {
	ticker := time.NewTicker(a.pollInterval)
	reportTicker := time.NewTicker(a.reportInterval)
	for {
		select {
		case <-ticker.C:
			a.poolCount++
			metrics := a.collectMetrics()
			metrics["PoolCount"] = float64(a.poolCount)
			fmt.Println("Metrics collected:", metrics)
		case <-reportTicker.C:
			metrics := a.collectMetrics()
			metrics["PoolCount"] = float64(a.poolCount)
			fmt.Println("Metrics collected:", metrics)

			for name, value := range metrics {
				a.sendMetrics("gauge", name, value)
			}
			a.sendMetrics("counter", "PoolCount", float64(a.poolCount))

		}
	}
}

func main() {
	agent := NewAgent(2*time.Second, 10*time.Second)
	agent.start()

}
