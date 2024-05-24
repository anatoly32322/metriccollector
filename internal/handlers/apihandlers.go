package apihandlers

import (
	"fmt"
	st "github.com/anatoly32322/metriccollector/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

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
		_, _ = w.Write([]byte(fmt.Sprint(memStorage.GetAll())))
	}
}
