package filemanager_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/memstorage"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/memstorage/filemanager"
	"github.com/stretchr/testify/assert"
)

func TestFileManager_SaveData(t *testing.T) {
	filePath := "./test_data.json"
	defer func() {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("Cant remove file: " + filePath)
		}
	}()

	fm := &filemanager.FileManager{}
	storage := &memstorage.MemStorage{
		Gauges:   map[string]float64{"testGauge": 42.42},
		Counters: map[string]int64{"testCounter": 100},
	}

	err := fm.SaveData(filePath, storage)
	assert.NoError(t, err)

	fileData, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.NotEmpty(t, fileData)

	var storageMap map[string]interface{}
	err = json.Unmarshal(fileData, &storageMap)
	assert.NoError(t, err)

	gauges := make(map[string]float64)
	for k, v := range storageMap["gauges"].(map[string]interface{}) {
		gauges[k] = v.(float64)
	}
	assert.Equal(t, storage.GetAllGauges(context.Background()), gauges)

	counters := make(map[string]int64)
	for k, v := range storageMap["counters"].(map[string]interface{}) {
		counters[k] = int64(v.(float64)) // Преобразуем float64 в int64
	}
	assert.Equal(t, storage.GetAllCounters(context.Background()), counters)
}

func TestFileManager_LoadData(t *testing.T) {
	filePath := "./test_data_load.json"
	defer func() {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("Cant remove file: " + filePath)
		}
	}()

	storageData := map[string]interface{}{
		"gauges":   map[string]float64{"testGauge": 42.42},
		"counters": map[string]int64{"testCounter": 100},
	}
	fileData, _ := json.Marshal(storageData)
	_ = os.WriteFile(filePath, fileData, 0644)

	fm := &filemanager.FileManager{}
	storage := &memstorage.MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	err := fm.LoadData(filePath, storage)
	assert.NoError(t, err)

	assert.Equal(t, float64(42.42), storage.Gauges["testGauge"])
	assert.Equal(t, int64(100), storage.Counters["testCounter"])
}

func TestFileManager_LoadEmptyData(t *testing.T) {
	filePath := "./test_empty_data.json"
	defer func() {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("Cant remove file: " + filePath)
		}
	}()

	_, _ = os.Create(filePath)

	fm := &filemanager.FileManager{}
	storage := &memstorage.MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	err := fm.LoadData(filePath, storage)
	assert.NoError(t, err)
	assert.Empty(t, storage.Gauges)
	assert.Empty(t, storage.Counters)
}

func TestFileManager_LoadNonExistentFile(t *testing.T) {
	filePath := "./non_existent_file.json"

	fm := &filemanager.FileManager{}
	storage := &memstorage.MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	err := fm.LoadData(filePath, storage)
	assert.NoError(t, err)
	assert.Empty(t, storage.Gauges)
	assert.Empty(t, storage.Counters)
}
