package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Host                  string `env:"ADDRESS"`
	PollIntervalSeconds   int64  `env:"POLL_INTERVAL"`
	ReportIntervalSeconds int64  `env:"REPORT_INTERVAL"`
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.Host, "a", "localhost:8080", "host")
	flag.Int64Var(&cfg.ReportIntervalSeconds, "r", 10, "report interval")
	flag.Int64Var(&cfg.PollIntervalSeconds, "p", 2, "poll interval")
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
	reportInterval := int64(time.Duration(cfg.ReportIntervalSeconds) * time.Second / pollInterval)
	log.Info(pollInterval)
	log.Info(reportInterval)
	var intervalCounter int64
	var pollCounter int64
	var gaugeMetrics map[string]float64
	for {
		if intervalCounter >= reportInterval {
			log.Info("sending metrics")
			gaugeMetrics = collectMetrics()
			client := resty.New()
			for metricName, metricValue := range gaugeMetrics {
				req, err := json.Marshal(&Metrics{
					ID:    metricName,
					MType: "gauge",
					Value: &metricValue,
				})
				if err != nil {
					log.Error(err)
				}
				_, err = client.R().
					SetHeader("Content-Type", "application/json").
					SetBody(req).
					Post(fmt.Sprintf("http://%s/update/", cfg.Host))
				if err != nil {
					log.Error(err)
				}
			}
			req, err := json.Marshal(&Metrics{
				ID:    "PollCount",
				MType: "counter",
				Delta: &pollCounter,
			})
			if err != nil {
				log.Error(err)
			}
			_, err = client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(req).
				Post(fmt.Sprintf("http://%s/update/", cfg.Host))
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
