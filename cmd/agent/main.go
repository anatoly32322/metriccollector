package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var (
	pollIntervalSeconds   int
	reportIntervalSeconds int
)

func main() {
	host := flag.String("host", "localhost:8080", "host")
	flag.IntVar(&reportIntervalSeconds, "r", 10, "report interval")
	flag.IntVar(&pollIntervalSeconds, "p", 2, "poll interval")
	flag.Parse()
	run(*host)
}

func run(host string) {
	pollInterval := time.Duration(pollIntervalSeconds) * time.Second
	reportInterval := int(time.Duration(reportIntervalSeconds) * time.Second / pollInterval)
	log.Info(pollInterval)
	log.Info(reportInterval)
	var intervalCounter = 0
	var pollCounter = 0
	var metrics map[string]string
	for {
		if intervalCounter >= reportInterval {
			log.Info("sending metrics")
			metrics = collectMetrics()
			pollCounter++
			client := resty.New()
			for metricName, metricValue := range metrics {
				_, err := client.R().SetPathParams(map[string]string{
					"metricType":  "gauge",
					"metricName":  metricName,
					"metricValue": metricValue,
				}).Post(fmt.Sprintf("http://%s/update/{metricType}/{metricName}/{metricValue}", host))
				if err != nil {
					log.Error(err)
				}
			}
			_, err := client.R().SetPathParams(map[string]string{
				"metricType":  "counter",
				"metricName":  "PollCount",
				"metricValue": strconv.Itoa(pollCounter),
			}).Post(fmt.Sprintf("http://%s/update/{metricType}/{metricName}/{metricValue}", host))
			if err != nil {
				log.Error(err)
			}
			intervalCounter = 0
		}
		time.Sleep(pollInterval)
		intervalCounter++
	}
}
