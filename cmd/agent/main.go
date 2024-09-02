package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/anatoly32322/metriccollector/internal/logger"
	"github.com/anatoly32322/metriccollector/internal/retry"
	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	Host                  string `env:"ADDRESS"`
	PollIntervalSeconds   int64  `env:"POLL_INTERVAL"`
	ReportIntervalSeconds int64  `env:"REPORT_INTERVAL"`
}

func main() {
	localLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer localLogger.Sync()

	logger.Sugar = *localLogger.Sugar()

	var cfg Config
	flag.StringVar(&cfg.Host, "a", "localhost:8080", "host")
	flag.Int64Var(&cfg.ReportIntervalSeconds, "r", 10, "report interval")
	flag.Int64Var(&cfg.PollIntervalSeconds, "p", 2, "poll interval")
	flag.Parse()

	err = env.Parse(&cfg)
	if err != nil {
		logger.Sugar.Fatal(err)
	}
	logger.Sugar.Info(cfg)

	run(cfg)
}

func run(cfg Config) {

	pollInterval := time.Duration(cfg.PollIntervalSeconds) * time.Second
	reportInterval := int64(time.Duration(cfg.ReportIntervalSeconds) * time.Second / pollInterval)
	logger.Sugar.Info(pollInterval)
	logger.Sugar.Info(reportInterval)
	var intervalCounter int64
	var pollCounter int64
	var gaugeMetrics map[string]float64
	var batch []Metrics
	for {
		if intervalCounter >= reportInterval {
			logger.Sugar.Info("sending metrics...")
			gaugeMetrics = collectMetrics()
			client := resty.New()
			for metricName, metricValue := range gaugeMetrics {
				batch = append(batch, Metrics{
					ID:    metricName,
					MType: "gauge",
					Value: &metricValue,
				})
			}
			batch = append(batch, Metrics{
				ID:    "PollCount",
				MType: "counter",
				Delta: &pollCounter,
			})
			req, err := json.Marshal(batch)
			if err != nil {
				logger.Sugar.Error(err)
			}
			send := func() error {
				_, err = client.R().
					SetHeader("Content-Type", "application/json").
					SetBody(req).
					Post(fmt.Sprintf("http://%s/updates/", cfg.Host))
				return err
			}
			_, err = retry.Retry(
				send,
				[]interface{}{},
				3, 1, 5,
			)

			if err != nil {
				logger.Sugar.Error(err)
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
