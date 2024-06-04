package storage

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Storage interface {
	Update(string, string, string) error
	UpdateV2(Metrics) (*Metrics, error)
	Get(string, string) (string, error)
	GetV2(Metrics) (*Metrics, error)
	GetAll() ([]byte, error)
}

type MemStorage struct {
	mx                 sync.Mutex
	GaugeMetrics       map[string]float64 `json:"gauge_metrics"`
	CounterMetrics     map[string]int64   `json:"counter_metrics"`
	AcceptedMetricType map[string]bool    `json:"-"`
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
		AcceptedMetricType: map[string]bool{
			"gauge":   true,
			"counter": true,
		},
	}
}

func (s *MemStorage) Update(metricType, metricName, value string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
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
		s.CounterMetrics[metricName] += intValue
	default:
		return fmt.Errorf("unknown metric type: %s", metricType)
	}
	return nil
}

func (s *MemStorage) UpdateV2(metric Metrics) (*Metrics, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if !s.AcceptedMetricType[metric.MType] {
		return nil, fmt.Errorf("metric type %s not accepted", metric.MType)
	}
	switch metric.MType {
	case "gauge":
		if metric.Value == nil {
			return nil, fmt.Errorf("metric value is nil")
		}
		s.GaugeMetrics[metric.ID] = *metric.Value
	case "counter":
		if metric.Delta == nil {
			return nil, fmt.Errorf("metric delta is nil")
		}
		s.CounterMetrics[metric.ID] += *metric.Delta
		*metric.Delta = s.CounterMetrics[metric.ID]
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metric.MType)
	}
	return &metric, nil
}

func (s *MemStorage) Get(metricType, metricName string) (string, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	switch metricType {
	case "gauge":
		if value, ok := s.GaugeMetrics[metricName]; ok {
			return fmt.Sprintf("%g", value), nil
		}
		return "", fmt.Errorf("gauge metric not found: %s", metricName)
	case "counter":
		if _, ok := s.CounterMetrics[metricName]; !ok {
			return "", fmt.Errorf("counter metric with name %s does not exist", metricName)
		}
		return strconv.FormatInt(s.CounterMetrics[metricName], 10), nil
	}
	return "", fmt.Errorf("unknown metric type: %s", metricType)
}

func (s *MemStorage) GetV2(metrics Metrics) (*Metrics, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	switch metrics.MType {
	case "gauge":
		if value, ok := s.GaugeMetrics[metrics.ID]; ok {
			metrics.Value = &value

			return &metrics, nil
		}

		return nil, fmt.Errorf("gauge metric not found: %s", metrics.ID)
	case "counter":
		if delta, ok := s.CounterMetrics[metrics.ID]; ok {
			metrics.Delta = &delta

			return &metrics, nil
		}

		return nil, fmt.Errorf("counter metric with name %s does not exist", metrics.ID)
	}

	return nil, fmt.Errorf("unknown metric type: %s", metrics.MType)
}

func (s *MemStorage) GetAll() ([]byte, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	return json.Marshal(s)
}
