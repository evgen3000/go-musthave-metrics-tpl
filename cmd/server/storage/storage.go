package storage

import (
	"fmt"
	"time"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/postgres"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/dbstorage"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/memstorage"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/memstorage/filemanager"
	"evgen3000/go-musthave-metrics-tpl.git/internal/dto"
)

type Interface interface {
	StorageType() string
	SetMetrics(dto []dto.MetricsDTO)
	SetGauge(metricName string, value float64)
	IncrementCounter(metricName string, value int64)
	GetGauge(metricName string) (float64, bool)
	GetCounter(metricName string) (int64, bool)
	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64
}

type Config struct {
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
	Database        string
}

func NewStorage(config Config, fm *filemanager.FileManager) (Interface, error) {
	if config.Database == "" {
		var storage Interface = &memstorage.MemStorage{
			Gauges:   make(map[string]float64),
			Counters: make(map[string]int64),
		}
		if config.Restore {
			err := fm.LoadData(config.FileStoragePath, storage)
			if err != nil {
				return nil, fmt.Errorf("can't read from storage file. %w", err)
			}
		}
		return storage, nil
	} else {
		pool := postgres.Connect(config.Database)
		var storage Interface = &dbstorage.DBStorage{Pool: pool}
		return storage, nil
	}

}
