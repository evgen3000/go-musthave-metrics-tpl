package filemanager

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type StorageInterface interface {
	SetGauge(metricName string, value float64)
	IncrementCounter(metricName string, value int64)
	GetGauge(metricName string) (float64, bool)
	GetCounter(metricName string) (int64, bool)
	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64
}

type FileManager struct{}

func (fm *FileManager) SaveData(filePath string, storage StorageInterface) error {
	storageMap := map[string]interface{}{
		"gauges":   storage.GetAllGauges(),
		"counters": storage.GetAllCounters(),
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
