package storage_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/memstorage"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/memstorage/filemanager"
	"github.com/stretchr/testify/assert"
)

func TestNewStorage_MemoryStorage(t *testing.T) {
	config := storage.Config{
		StoreInterval:   time.Second,
		FileStoragePath: "./test_storage.json",
		Restore:         true,
		Database:        "",
	}
	fm := &filemanager.FileManager{}

	storageInstance, err := storage.NewStorage(config, fm)
	assert.NoError(t, err)
	assert.NotNil(t, storageInstance)
	assert.Equal(t, "ms", storageInstance.StorageType())
}

func TestNewStorage_MemoryStorageWithRestore(t *testing.T) {
	filePath := "./empty_test_data.json"
	defer func() {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("Cant remove file: " + filePath)
		}
	}()

	_, err := os.Create(filePath)
	assert.NoError(t, err, "Failed to create empty test data file.")

	fm := &filemanager.FileManager{}

	config := storage.Config{
		FileStoragePath: filePath,
		Restore:         true,
	}

	storageInstance, err := storage.NewStorage(config, fm)
	assert.NoError(t, err, "Expected no error when creating storage instance with empty restore file.")

	assert.NotNil(t, storageInstance, "Storage instance should not be nil.")

	memStorage, ok := storageInstance.(*memstorage.MemStorage)
	assert.True(t, ok, "Storage instance should be of type *memstorage.MemStorage.")
	assert.Empty(t, memStorage.Gauges, "Gauges should be empty after restoration from an empty file.")
	assert.Empty(t, memStorage.Counters, "Counters should be empty after restoration from an empty file.")
}
