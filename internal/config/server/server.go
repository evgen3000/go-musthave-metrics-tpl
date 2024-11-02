package server

import (
	"flag"
	"time"

	"evgen3000/go-musthave-metrics-tpl.git/internal/config/utils"
)

type Config struct {
	Host          string
	StoreInterval time.Duration
	FilePath      string
	Restore       bool
	Database      string
}

func GetServerConfig() *Config {
	hostFlag := flag.String("a", "localhost:8080", "Host IP address and port.")
	storeIntervalFlag := flag.Int("i", 300, "Store interval in sec.")
	filePathFlag := flag.String("f", "storage.json", "File storage location.")
	restoreFlag := flag.Bool("r", true, "Restore stored configuration.")
	databaseURL := flag.String("d", "postgres://admin:admin@localhost:5432/admin", "Database IP address and port. like: postgres://admin:admin@localhost:5432/admin")
	flag.Parse()
	return &Config{
		Host:          utils.GetStringValue("ADDRESS", *hostFlag),
		FilePath:      utils.GetStringValue("FILE_STORE_PATH", *filePathFlag),
		StoreInterval: time.Duration(utils.GetIntValue("STORE_INTERVAL", *storeIntervalFlag)) * time.Second,
		Restore:       utils.GetBoolValue("RESTORE", *restoreFlag),
		Database:      utils.GetStringValue("DATABASE_DSN", *databaseURL),
	}
}
