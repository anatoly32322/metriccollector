package main

import (
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var (
	pollInterval   = 2 * time.Second
	reportInterval = int(10 * time.Second / pollInterval)
)

func main() {
	run()
}

func run() {
	var intervalCounter = 0
	var pollCounter = 0
	var metrics map[string]string
	for {
		if intervalCounter >= reportInterval {
			metrics = collectMetrics()
			pollCounter++
			client := resty.New()
			for metricName, metricValue := range metrics {
				_, err := client.R().SetPathParams(map[string]string{
					"metricType":  "gauge",
					"metricName":  metricName,
					"metricValue": metricValue,
				}).Post("http://localhost:8080/update/{metricType}/{metricName}/{metricValue}")
				if err != nil {
					log.Error(err)
				}
			}
			_, err := client.R().SetPathParams(map[string]string{
				"metricType":  "counter",
				"metricName":  "PollCount",
				"metricValue": strconv.Itoa(pollCounter),
			}).Post("http://localhost:8080/update/{metricType}/{metricName}/{metricValue}")
			if err != nil {
				log.Error(err)
			}
			intervalCounter = 0
		}
		time.Sleep(pollInterval)
		intervalCounter++
	}
}
