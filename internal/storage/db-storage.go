package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anatoly32322/metriccollector/internal/logger"
	"strconv"
	"sync"
)

var (
	createTables = `
		CREATE SCHEMA IF NOT EXISTS metric_collector;
		CREATE TABLE IF NOT EXISTS metric_collector.gauge_metrics (
		    metric_name varchar(32) UNIQUE,
		    value double precision
		);
		CREATE TABLE IF NOT EXISTS metric_collector.counter_metrics (
		    metric_name varchar(32) UNIQUE,
		    delta bigint
		);
	`
	insertGaugeQuery = `
		INSERT INTO metric_collector.gauge_metrics (metric_name, value)
		VALUES ($1, $2)
		ON CONFLICT (metric_name) 
		DO UPDATE SET 
		    value=EXCLUDED.value;
	`
	insertCounterQuery = `
		INSERT INTO metric_collector.counter_metrics (metric_name, delta)
		VALUES ($1, $2)
		ON CONFLICT (metric_name)
		DO UPDATE SET
		    delta=counter_metrics.delta + EXCLUDED.delta;
	`
	selectGaugeQuery = `
		SELECT value FROM metric_collector.gauge_metrics
		WHERE metric_name = $1;
	`
	selectCounterQuery = `
		SELECT delta FROM metric_collector.counter_metrics
		WHERE metric_name = $1;
    `
	selectAllGaugeQuery = `
		SELECT metric_name, value FROM metric_collector.gauge_metrics;
	`
	selectAllCounterQuery = `
		SELECT metric_name, delta FROM metric_collector.counter_metrics;
	`
)

type DBStorage struct {
	mx *sync.Mutex
	db *sql.DB
}

func NewDBStorage(dsn string) (*DBStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = initDB(db)
	if err != nil {
		return nil, err
	}
	return &DBStorage{mx: &sync.Mutex{}, db: db}, nil
}

func initDB(db *sql.DB) error {
	_, err := db.Exec(createTables)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBStorage) Update(metricType, metricName, value string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	logger.Sugar.Infof("got metric: %s, %s, %s", metricType, metricName, value)
	switch metricType {
	case "gauge":
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logger.Sugar.Errorf("got error: %e", err)
			return err
		}
		logger.Sugar.Infof("exec update query with args: %s, %f", metricName, floatValue)
		_, err = s.db.Exec(insertGaugeQuery, metricName, floatValue)
		if err != nil {
			logger.Sugar.Errorf("got error: %e", err)
			return fmt.Errorf("got error during exec update query: %e", err)
		}
	case "counter":
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			logger.Sugar.Errorf("got error: %e", err)
			return err
		}
		_, err = s.db.Exec(insertCounterQuery, metricName, intValue)
		if err != nil {
			logger.Sugar.Errorf("got error: %e", err)
			return fmt.Errorf("got error during exec update query: %e", err)
		}
	default:
		return fmt.Errorf("unknown metric type: %s", metricType)
	}
	return nil
}

func (s *DBStorage) UpdateV2(metric Metric) (err error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	logger.Sugar.Infof("got metric: %s, %s", metric.MType, metric.ID)
	switch metric.MType {
	case "gauge":
		if metric.Value == nil {
			err = fmt.Errorf("metric value is nil")
			logger.Sugar.Errorf("got error: %e", err)
			return
		}
		logger.Sugar.Infof("exec update v2 query with args: %s, %d", metric.ID, metric.Value)
		_, err = s.db.Exec(insertGaugeQuery, metric.ID, metric.Value)
		if err != nil {
			logger.Sugar.Errorf("got error: %e", err)
			return fmt.Errorf("got error during exec update v2 query: %e", err)
		}
		return
	case "counter":
		if metric.Delta == nil {
			err = fmt.Errorf("metric delta is nil")
			logger.Sugar.Errorf("got error: %e", err)
			return
		}
		_, err = s.db.Exec(insertCounterQuery, metric.ID, metric.Delta)
		if err != nil {
			logger.Sugar.Errorf("got error: %e", err)
			return fmt.Errorf("got error during exec update v2 query: %e", err)
		}
		return
	default:
		return fmt.Errorf("unknown metric type: %s", metric.MType)
	}
}

func (s *DBStorage) UpdateBatch(metrics []Metric) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	logger.Sugar.Infof("got metrics: %v", metrics)
	tx, err := s.db.Begin()
	if err != nil {
		logger.Sugar.Errorf("failed init transaction: %s", err)
		return fmt.Errorf("failed init transaction: %w", err)
	}
	for i := 0; i < len(metrics); i++ {
		switch metrics[i].MType {
		case "gauge":
			if metrics[i].Value == nil {
				err = fmt.Errorf("metric value is nil")
				logger.Sugar.Errorf("got error: %e", err)
				tx.Rollback()
				return err
			}
			logger.Sugar.Infof("exec update batch query with args: %s, %d", metrics[i].ID, metrics[i].Value)
			_, err = tx.Exec(insertGaugeQuery, metrics[i].ID, metrics[i].Value)
			if err != nil {
				logger.Sugar.Errorf("got error: %e", err)
				tx.Rollback()
				return fmt.Errorf("got error during exec update batch query: %e", err)
			}
		case "counter":
			if metrics[i].Delta == nil {
				err = fmt.Errorf("metric delta is nil")
				logger.Sugar.Errorf("got error: %e", err)
				tx.Rollback()
				return err
			}
			logger.Sugar.Infof("exec update batch query with args: %s, %d", metrics[i].ID, metrics[i].Value)
			_, err = s.db.Exec(insertCounterQuery, metrics[i].ID, metrics[i].Delta)
			if err != nil {
				logger.Sugar.Errorf("got error: %e", err)
				tx.Rollback()
				return fmt.Errorf("got error during exec update batch query: %e", err)
			}
		default:
			tx.Rollback()
			return fmt.Errorf("unknown metric type: %s", metrics[i].MType)
		}
	}
	return tx.Commit()
}

func (s *DBStorage) Get(metricType, metricName string) (res string, err error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	switch metricType {
	case "gauge":
		logger.Sugar.Infof("exec get query with arg: %s", metricName)
		row := s.db.QueryRow(selectGaugeQuery, metricName)

		if err = row.Scan(&res); err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Sugar.Errorf("got error: %e", err)
			return
		}
		return
	case "counter":
		logger.Sugar.Infof("exec get query with arg: %s", metricName)
		row := s.db.QueryRow(selectCounterQuery, metricName)

		if err = row.Scan(&res); err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Sugar.Errorf("got error: %e", err)
			return
		}
		return
	}
	return "", fmt.Errorf("unknown metric type: %s", metricType)
}

func (s *DBStorage) GetV2(metric Metric) (res *Metric, err error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	res = &metric

	switch metric.MType {
	case "gauge":
		logger.Sugar.Infof("exec get v2 query with arg: %s", metric.ID)
		row := s.db.QueryRow(selectGaugeQuery, metric.ID)

		var val float64
		if err = row.Scan(&val); err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Sugar.Errorf("got error: %e", err)
			return
		}
		res.Value = &val
		return
	case "counter":
		logger.Sugar.Infof("exec get v2 query with arg: %s", metric.ID)
		row := s.db.QueryRow(selectCounterQuery, metric.ID)

		var val int64
		if err = row.Scan(&val); err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Sugar.Errorf("got error: %e", err)
			return
		}
		res.Delta = &val
		return
	}
	return res, fmt.Errorf("unknown metric type: %s", metric.MType)
}

func (s *DBStorage) GetAll() ([]byte, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	var metric Metric
	type resultStorage struct {
		GaugeMetrics   map[string]float64 `json:"gauge_metrics"`
		CounterMetrics map[string]int64   `json:"counter_metrics"`
	}
	results := &resultStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
	}
	logger.Sugar.Info("exec get all query")
	rows, err := s.db.Query(selectAllGaugeQuery)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed exec get all query: %w", err)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed exec get all query: %w", err)
	}
	for rows.Next() {
		err = rows.Scan(&metric.ID, &metric.Value)
		if err != nil {
			return nil, fmt.Errorf("failed scan row: %w", err)
		}
		results.GaugeMetrics[metric.ID] = *metric.Value
	}

	rows, err = s.db.Query(selectAllCounterQuery)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed exec get all query: %w", err)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed exec get all query: %w", err)
	}
	for rows.Next() {
		err = rows.Scan(&metric.ID, &metric.Delta)
		if err != nil {
			return nil, fmt.Errorf("failed scan row: %w", err)
		}
		results.CounterMetrics[metric.ID] = *metric.Delta
	}
	return json.Marshal(results)
}

func (s *DBStorage) Save(fname string) error {
	panic("save not implemented")
}

func (s *DBStorage) Load(fname string) error {
	panic("load not implemented")
}
