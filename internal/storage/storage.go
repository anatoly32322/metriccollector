package storage

import (
	"fmt"
	"strconv"
)

type MemStorage struct {
	GaugeMetrics       map[string]float64
	CounterMetrics     map[string][]int64
	AcceptedMetricType map[string]bool
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string][]int64),
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
		s.GaugeMetrics[metricName] = floatValue
	case "counter":
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		s.CounterMetrics[metricName] = append(s.CounterMetrics[metricName], intValue)
	default:
		return fmt.Errorf("unknown metric type: %s", metricType)
	}
	return nil
}

func (s *MemStorage) Get(metricType, metricName string) (string, error) {
	switch metricType {
	case "gauge":
		if value, ok := s.GaugeMetrics[metricName]; ok {
			return strconv.FormatFloat(value, 'f', 6, 64), nil
		}
		return "", fmt.Errorf("gauge metric not found: %s", metricName)
	case "counter":
		if len(s.CounterMetrics[metricName]) == 0 {
			return "", fmt.Errorf("counter metric with name %s does not exist", metricName)
		}
		return strconv.FormatInt(s.CounterMetrics[metricName][len(s.CounterMetrics[metricName])-1], 64), nil
	}
	return "", fmt.Errorf("unknown metric type: %s", metricType)
}
