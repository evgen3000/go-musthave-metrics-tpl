package memstorage_test

import (
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
	s := storage.NewStorage(StorageConfig, &fm)
	s.SetGauge("temperature", 23.5)

	value, exists := s.GetGauge("temperature")
	assert.True(t, exists)
	assert.Equal(t, 23.5, value)
}
func TestMemStorage_IncrementCounter(t *testing.T) {
	fm := filemanager.FileManager{}
	s := storage.NewStorage(StorageConfig, &fm)
	s.IncrementCounter("hits", 10)
	s.IncrementCounter("hits", 5)

	value, exists := s.GetCounter("hits")
	assert.True(t, exists)
	assert.Equal(t, int64(15), value)
}
