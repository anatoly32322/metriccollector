package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Config struct {
	Host                  string `env:"ADDRESS"`
	PollIntervalSeconds   int    `env:"POLL_INTERVAL"`
	ReportIntervalSeconds int    `env:"REPORT_INTERVAL"`
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.Host, "a", "localhost:8080", "host")
	flag.IntVar(&cfg.ReportIntervalSeconds, "r", 10, "report interval")
	flag.IntVar(&cfg.PollIntervalSeconds, "p", 2, "poll interval")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg)

	run(cfg)
}

func run(cfg Config) {
	pollInterval := time.Duration(cfg.PollIntervalSeconds) * time.Second
	reportInterval := int(time.Duration(cfg.ReportIntervalSeconds) * time.Second / pollInterval)
	log.Info(pollInterval)
	log.Info(reportInterval)
	var intervalCounter = 0
	var pollCounter = 0
	var metrics map[string]string
	for {
		if intervalCounter >= reportInterval {
			log.Info("sending metrics")
			metrics = collectMetrics()
			client := resty.New()
			for metricName, metricValue := range metrics {
				_, err := client.R().SetPathParams(map[string]string{
					"metricType":  "gauge",
					"metricName":  metricName,
					"metricValue": metricValue,
				}).Post(fmt.Sprintf("http://%s/update/{metricType}/{metricName}/{metricValue}", cfg.Host))
				if err != nil {
					log.Error(err)
				}
			}
			_, err := client.R().SetPathParams(map[string]string{
				"metricType":  "counter",
				"metricName":  "PollCount",
				"metricValue": strconv.Itoa(pollCounter),
			}).Post(fmt.Sprintf("http://%s/update/{metricType}/{metricName}/{metricValue}", cfg.Host))
			if err != nil {
				log.Error(err)
			} else {
				pollCounter = 0
			}
			intervalCounter = 0
		}
		pollCounter++
		time.Sleep(pollInterval)
		intervalCounter++
	}
}
