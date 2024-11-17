package agent

import (
	"flag"

	"evgen3000/go-musthave-metrics-tpl.git/internal/config/utils"
)

type Config struct {
	PoolInterval   int
	ReportInterval int
	Host           string
	CryptoKey      string
}

func GetAgentConfig() *Config {
	reportIntervalFlag := flag.Int("r", 10, "Report interval in seconds.")
	pollIntervalFlag := flag.Int("p", 2, "Pool interval in seconds.")
	hostFlag := flag.String("a", "localhost:8080", "Host IP address and port.")
	cryptoKey := flag.String("k", "", "AES encryption key.")
	flag.Parse()

	return &Config{
		PoolInterval:   utils.GetIntValue("POLL_INTERVAL", *pollIntervalFlag),
		ReportInterval: utils.GetIntValue("REPORT_INTERVAL", *reportIntervalFlag),
		Host:           utils.GetStringValue("ADDRESS", *hostFlag),
		CryptoKey:      *cryptoKey,
	}
}
