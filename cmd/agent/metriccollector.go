package main

import (
	"math/rand"
	"runtime"
	"strconv"
)

func collectMetrics() map[string]string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	metrics := map[string]string{}
	metrics["Alloc"] = strconv.FormatUint(m.Alloc, 10)
	metrics["BuckHashSys"] = strconv.FormatUint(m.BuckHashSys, 10)
	metrics["Frees"] = strconv.FormatUint(m.Frees, 10)
	metrics["GCCPUFraction"] = strconv.FormatFloat(m.GCCPUFraction, 'f', 3, 64)
	metrics["GCSys"] = strconv.FormatUint(m.GCSys, 10)
	metrics["HeapAlloc"] = strconv.FormatUint(m.HeapAlloc, 10)
	metrics["HeapIdle"] = strconv.FormatUint(m.HeapIdle, 10)
	metrics["HeapInuse"] = strconv.FormatUint(m.HeapInuse, 10)
	metrics["HeapReleased"] = strconv.FormatUint(m.HeapReleased, 10)
	metrics["HeapObjects"] = strconv.FormatUint(m.HeapObjects, 10)
	metrics["HeapSys"] = strconv.FormatUint(m.HeapSys, 10)
	metrics["LastGC"] = strconv.FormatUint(m.LastGC, 10)
	metrics["Lookups"] = strconv.FormatUint(m.Lookups, 10)
	metrics["MCacheInuse"] = strconv.FormatUint(m.MCacheInuse, 10)
	metrics["MCacheSys"] = strconv.FormatUint(m.MCacheSys, 10)
	metrics["MSpanInuse"] = strconv.FormatUint(m.MSpanInuse, 10)
	metrics["MSpanSys"] = strconv.FormatUint(m.MSpanSys, 10)
	metrics["Mallocs"] = strconv.FormatUint(m.Mallocs, 10)
	metrics["NextGC"] = strconv.FormatUint(m.NextGC, 10)
	metrics["NumForcedGC"] = strconv.FormatUint(uint64(m.NumForcedGC), 10)
	metrics["NumGC"] = strconv.FormatUint(uint64(m.NumGC), 10)
	metrics["OtherSys"] = strconv.FormatUint(m.OtherSys, 10)
	metrics["PauseTotalNs"] = strconv.FormatUint(m.PauseTotalNs, 10)
	metrics["StackInuse"] = strconv.FormatUint(m.StackInuse, 10)
	metrics["StackSys"] = strconv.FormatUint(m.StackSys, 10)
	metrics["Sys"] = strconv.FormatUint(m.Sys, 10)
	metrics["TotalAlloc"] = strconv.FormatUint(m.TotalAlloc, 10)

	metrics["GCCPUFraction"] = strconv.FormatFloat(rand.Float64(), 'f', 3, 64)
	return metrics
}
