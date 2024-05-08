package apihandlers

import (
	st "github.com/anatoly32322/metriccollector/internal/storage"
	"github.com/go-chi/chi/v5"
)

func MetricRouter(memStorage *st.MemStorage) chi.Router {
	router := chi.NewRouter()

	router.Post("/update/{metricType}/{metricName}/{metricValue}", ServeUpdateHandler(memStorage))
	router.Get("/value/{metricType}/{metricName}", GetMetricHandler(memStorage))

	return router
}
