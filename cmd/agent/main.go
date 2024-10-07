package main

import (
	"evgen3000/go-musthave-metrics-tpl.git/cmd/agent/client"
	"flag"
	"os"
	"strconv"
	"time"
)

func main() {
	hostEnv, isHostEnv := os.LookupEnv("ADDRESS")
	reportIntervalEnv, isReportIntervalEn := os.LookupEnv("REPORT_INTERVAL")
	pollIntervalEnv, isPollIntervalEnv := os.LookupEnv("POLL_INTERVAL")

	hostFlag := flag.String("a", "localhost:8080", "Host IP address and port.")
	reportIntervalFlag := flag.Int("r", 10, "Report interval in seconds.")
	pollIntervalFlag := flag.Int("p", 2, "Pool interval in seconds.")
	flag.Parse()

	if isReportIntervalEn && isPollIntervalEnv && isHostEnv {
		poolInterval, err := strconv.ParseInt(pollIntervalEnv, 10, 64)
		if err != nil {
			panic(err)
		}
		reportInterval, err := strconv.ParseInt(reportIntervalEnv, 10, 64)
		if err != nil {
			panic(err)
		}
		(client.NewAgent(hostEnv, time.Duration(poolInterval)*time.Second, time.Duration((reportInterval))*time.Second)).Start()
	} else {
		(client.NewAgent(*hostFlag, time.Duration(*pollIntervalFlag)*time.Second, time.Duration((*reportIntervalFlag))*time.Second)).Start()
	}

}
