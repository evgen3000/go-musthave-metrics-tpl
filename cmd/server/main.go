package main

import (
	"log"
	"net/http"
	"time"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/router"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/memstorage/filemanager"
	"evgen3000/go-musthave-metrics-tpl.git/internal/config/server"
	"evgen3000/go-musthave-metrics-tpl.git/internal/workerpool"
)

func main() {
	conf := server.GetServerConfig()
	fm := filemanager.FileManager{}
	s, err := storage.NewStorage(storage.Config{
		StoreInterval:   conf.StoreInterval,
		FileStoragePath: conf.FilePath,
		Restore:         conf.Restore,
		Database:        conf.Database,
	}, &fm)
	if err != nil {
		log.Fatal("Error initializing storage:", err)
	}

	// Инициализация пула воркеров
	workerPool := workerpool.NewWorkerPool(5, 100)

	r := router.SetupRouter(s, conf.CryptoKey, workerPool)

	ticker := time.NewTicker(conf.StoreInterval)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			if err := fm.SaveData(conf.FilePath, s); err != nil {
				log.Fatal("Can't save data to storage.json")
			} else {
				log.Println("Saved data")
			}
		}
	}()

	log.Println("Server is running on", conf.Host)
	log.Fatal(http.ListenAndServe(conf.Host, r))
}
