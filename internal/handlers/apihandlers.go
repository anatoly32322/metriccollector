package apihandlers

import (
	"fmt"
	st "github.com/anatoly32322/metriccollector/internal/storage"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func ServeUpdateHandler(memStorage *st.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// /update/{metricType}/{metricName}/{metricValue}

		log.Info(fmt.Sprintf("got request with path: %s", r.URL.Path))

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		if metricName == "" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(fmt.Sprintf("metric name not specified")))
			log.Error(fmt.Sprintf("metric name not specified: %s", metricName))
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

func GetMetricHandler(memStorage *st.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")

		value, err := memStorage.Get(metricType, metricName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(value))
	}
}
