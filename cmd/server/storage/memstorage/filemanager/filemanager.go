package filemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"evgen3000/go-musthave-metrics-tpl.git/internal/dto"
)

type StorageInterface interface {
	SetMetrics(ctx context.Context, dto []dto.MetricsDTO)
	SetGauge(ctx context.Context, metricName string, value float64)
	IncrementCounter(ctx context.Context, metricName string, value int64)
	GetGauge(ctx context.Context, metricName string) (float64, bool)
	GetCounter(ctx context.Context, metricName string) (int64, bool)
	GetAllGauges(ctx context.Context) map[string]float64
	GetAllCounters(ctx context.Context) map[string]int64
}

type FileManager struct{}

func (fm *FileManager) SaveData(filePath string, storage StorageInterface) error {
	storageMap := map[string]interface{}{
		"gauges":   storage.GetAllGauges(context.Background()),
		"counters": storage.GetAllCounters(context.Background()),
	}
	fileData, err := json.Marshal(storageMap)
	if err != nil {
		return fmt.Errorf("error while marshalling file: %w", err)
	}
	return os.WriteFile(filePath, fileData, 0644)
}

func (fm *FileManager) LoadData(filePath string, storage StorageInterface) error {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		return fmt.Errorf("can't open file. %w", err)
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("can't read file. %w", err)
	}
	if len(fileData) == 0 {
		log.Println("Storage file is empty, nothing to load.")
		return nil
	}

	err = json.Unmarshal(fileData, &storage)
	if err != nil {
		return fmt.Errorf("can't read json. %w", err)
	}
	return nil
}
