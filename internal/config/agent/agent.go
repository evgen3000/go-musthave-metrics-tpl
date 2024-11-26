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
	RateLimit      int
}

func GetAgentConfig() *Config {
	reportIntervalFlag := flag.Int("r", 10, "Report interval in seconds.")
	pollIntervalFlag := flag.Int("p", 2, "Pool interval in seconds.")
	rateLimit := flag.Int("l", 2, "Rate limit in seconds.")
	hostFlag := flag.String("a", "localhost:8080", "Host IP address and port.")
	cryptoKey := flag.String("k", "", "AES encryption key.")

	flag.Parse()

	return &Config{
		PoolInterval:   utils.GetIntValue("POLL_INTERVAL", *pollIntervalFlag),
		ReportInterval: utils.GetIntValue("REPORT_INTERVAL", *reportIntervalFlag),
		Host:           utils.GetStringValue("ADDRESS", *hostFlag),
		CryptoKey:      utils.GetStringValue("KEY", *cryptoKey),
		RateLimit:      utils.GetIntValue("RATE_LIMIT", *rateLimit),
	}
}
