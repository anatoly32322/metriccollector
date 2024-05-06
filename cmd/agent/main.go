package main

import (
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
	for {
		if intervalCounter >= reportInterval {
			_ = collectMetrics()
			pollCounter++
			// Отправка метрик через httpclient
			intervalCounter = 0
		}
		time.Sleep(pollInterval)
		intervalCounter++
	}
}
