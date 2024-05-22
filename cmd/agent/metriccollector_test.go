package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	collectableMetrics = []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapReleased",
		"HeapObjects",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
	}
)

func TestCheckCollectedMetrics(t *testing.T) {
	t.Run("check collected metrics contains every required metrics", func(t *testing.T) {
		metrics := collectMetrics()
		for _, metric := range collectableMetrics {
			assert.Contains(t, metrics, metric)
		}
	})
}
