package apihandlers

import (
	"bytes"
	"encoding/json"
	st "github.com/anatoly32322/metriccollector/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func ServeUpdateHandlerV2(memStorage st.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics st.Metrics
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

		computedMetrics, err := memStorage.UpdateV2(metrics)

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

func ServeUpdateHandler(memStorage st.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		if metricName == "" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("metric name not specified"))
			return
		}
		metricValue := chi.URLParam(r, "metricValue")

		err := memStorage.Update(metricType, metricName, metricValue)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))

			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func GetMetricHandlerV2(memStorage st.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics st.Metrics
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

		gotMetric, err := memStorage.GetV2(metrics)
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

func GetMetricHandler(memStorage st.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		value, err := memStorage.Get(metricType, metricName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(err.Error()))

			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(value))
	}
}

func GetPageHandler(memStorage st.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := memStorage.GetAll()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(data)
	}
}
