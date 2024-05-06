package storage

import (
	"fmt"
	"strconv"
)

type GaugeMetric struct {
	Name  string
	Value float64
}

type CounterMetric struct {
	Name  string
	Value int64
}

type MemStorage struct {
	GaugeMetrics       GaugeMetric
	CounterMetrics     []CounterMetric
	AcceptedMetricType map[string]bool
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMetrics:   GaugeMetric{},
		CounterMetrics: make([]CounterMetric, 0),
		AcceptedMetricType: map[string]bool{
			"gauge":   true,
			"counter": true,
		},
	}
}

func (s *MemStorage) Update(metricType, metricName, value string) error {
	if !s.AcceptedMetricType[metricType] {
		return fmt.Errorf("metric type %s not accepted", metricType)
	}
	switch metricType {
	case "gauge":
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		s.GaugeMetrics = GaugeMetric{metricName, floatValue}
	case "counter":
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		s.CounterMetrics = append(s.CounterMetrics, CounterMetric{metricName, intValue})
	default:
		return fmt.Errorf("unknown metric type: %s", metricType)
	}
	return nil
}
