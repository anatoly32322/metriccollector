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
	gaugeMetrics       GaugeMetric
	counterMetrics     []CounterMetric
	acceptedMetricType map[string]bool
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gaugeMetrics:   GaugeMetric{},
		counterMetrics: make([]CounterMetric, 0),
		acceptedMetricType: map[string]bool{
			"gauge":   true,
			"counter": true,
		},
	}
}

func (s *MemStorage) Update(metricType, metricName, value string) error {
	if !s.acceptedMetricType[metricType] {
		return fmt.Errorf("metric type %s not accepted", metricType)
	}
	switch metricType {
	case "gauge":
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		s.gaugeMetrics = GaugeMetric{metricName, floatValue}
	case "counter":
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		s.counterMetrics = append(s.counterMetrics, CounterMetric{metricName, intValue})
	default:
		return fmt.Errorf("unknown metric type: %s", metricType)
	}
	return nil
}
