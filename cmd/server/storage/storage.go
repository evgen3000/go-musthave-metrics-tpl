package storage

import (
	"context"
	"fmt"
	"time"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/postgres"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/dbstorage"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/memstorage"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/memstorage/filemanager"
	"evgen3000/go-musthave-metrics-tpl.git/internal/dto"
)

type MetricsStorage interface {
	StorageType() string
	SetMetrics(ctx context.Context, dto []dto.MetricsDTO)
	SetGauge(ctx context.Context, metricName string, value float64)
	IncrementCounter(ctx context.Context, metricName string, value int64)
	GetGauge(ctx context.Context, metricName string) (float64, bool)
	GetCounter(ctx context.Context, metricName string) (int64, bool)
	GetAllGauges(ctx context.Context) map[string]float64
	GetAllCounters(ctx context.Context) map[string]int64
}

type Config struct {
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
	Database        string
}

func NewStorage(config Config, fm *filemanager.FileManager) (MetricsStorage, error) {
	if config.Database == "" {
		var storage MetricsStorage = &memstorage.MemStorage{
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
		var storage MetricsStorage = &dbstorage.DBStorage{Pool: pool}
		return storage, nil
	}

}
