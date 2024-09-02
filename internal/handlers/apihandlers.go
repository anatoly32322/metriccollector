package apihandlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	log "github.com/anatoly32322/metriccollector/internal/logger"
	st "github.com/anatoly32322/metriccollector/internal/storage"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func ServeUpdatesHandler(storage st.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics []st.Metric
		var buf bytes.Buffer

		if r.Header.Get("Content-Type") != "application/json" {
			log.Sugar.Error("request body is not JSON")
			http.Error(w, "Request body is not JSON", http.StatusNotFound)
			return
		}

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			log.Sugar.Errorf("failed reading body: %s", err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
			log.Sugar.Errorf("failed unmarshal bode: %s", err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = storage.UpdateBatch(metrics)
		if err != nil {
			log.Sugar.Errorf("failed update batch: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func ServeUpdateHandlerV2(storage st.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics st.Metric
		var buf bytes.Buffer

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Request body is not JSON", http.StatusBadRequest)

			return
		}
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}
		if metrics.ID == "" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("metric name not specified"))
			return
		}

		err = storage.UpdateV2(metrics)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		computedMetrics, err := storage.GetV2(metrics)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		resp, err := json.Marshal(computedMetrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	}
}

func ServeUpdateHandler(storage st.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		if metricName == "" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("metric name not specified"))
			return
		}
		metricValue := chi.URLParam(r, "metricValue")

		err := storage.Update(metricType, metricName, metricValue)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))

			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func GetMetricHandlerV2(storage st.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics st.Metric
		var buf bytes.Buffer

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Request body is not JSON", http.StatusNotFound)

			return
		}
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &metrics); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if metrics.ID == "" {
			http.Error(w, "metric name not specified", http.StatusNotFound)
			return
		}

		gotMetric, err := storage.GetV2(metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		resp, err := json.Marshal(gotMetric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	}
}

func GetMetricHandler(storage st.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		value, err := storage.Get(metricType, metricName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(value))
	}
}

func GetPageHandler(storage st.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := storage.GetAll()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(data)
	}
}

func GetPing(dsn string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("pgx", dsn)
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Sugar.Errorf("failed close db connection: %e", err)
			}
		}(db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		err = db.Ping()
		if err != nil {
			log.Sugar.Warn("no database connection")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
