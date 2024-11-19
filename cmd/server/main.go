package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/router"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage"
	"evgen3000/go-musthave-metrics-tpl.git/cmd/server/storage/memstorage/filemanager"
	"evgen3000/go-musthave-metrics-tpl.git/internal/config/server"
	httpLogger "evgen3000/go-musthave-metrics-tpl.git/internal/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func runServer(config *server.Config, router *chi.Mux) {
	logger := httpLogger.InitLogger()
	logger.Info("server is running on", zap.String("host", config.Host))

	err := http.ListenAndServe(config.Host, router)
	if err != nil {
		logger.Fatal("Error", zap.String("Error", err.Error()))
	}
}

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
		log.Fatal(errors.Unwrap(err))
	}
	r := router.SetupRouter(s, conf.CryptoKey)

	ticker := time.NewTicker(conf.StoreInterval)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			if err := fm.SaveData(conf.FilePath, s); err != nil {
				log.Fatal("Can't to save data to storage.json")
			} else {
				log.Println("Saved data")
			}
		}
	}()

	runServer(conf, r)
}
