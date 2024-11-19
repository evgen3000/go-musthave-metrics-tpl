package memstorage_test

import (
	"context"
	"testing"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/memstorage/filemanager"
	"github.com/stretchr/testify/assert"
)

var StorageConfig = storage.Config{
	StoreInterval:   300,
	FileStoragePath: "storage.json",
	Restore:         true,
}

func TestMemStorageSetAndGetGauge(t *testing.T) {
	fm := filemanager.FileManager{}
	s, err := storage.NewStorage(StorageConfig, &fm)
	if err != nil {
		panic(err)
	}
	s.SetGauge(context.Background(), "temperature", 23.5)

	value, exists := s.GetGauge(context.Background(), "temperature")
	assert.True(t, exists)
	assert.Equal(t, 23.5, value)
}
func TestMemStorage_IncrementCounter(t *testing.T) {
	fm := filemanager.FileManager{}
	s, err := storage.NewStorage(StorageConfig, &fm)
	if err != nil {
		panic(err)
	}
	s.IncrementCounter(context.Background(), "hits", 10)
	s.IncrementCounter(context.Background(), "hits", 5)

	value, exists := s.GetCounter(context.Background(), "hits")
	assert.True(t, exists)
	assert.Equal(t, int64(15), value)
}
