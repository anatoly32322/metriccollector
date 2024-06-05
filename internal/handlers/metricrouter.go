package apihandlers

import (
	st "github.com/anatoly32322/metriccollector/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MetricRouter(memStorage st.Storage) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.Compress(5))

	router.Post("/update/{metricType}/{metricName}/{metricValue}", ServeUpdateHandler(memStorage))
	router.Post("/update/", ServeUpdateHandlerV2(memStorage))
	router.Post("/value/", GetMetricHandlerV2(memStorage))
	router.Get("/value/{metricType}/{metricName}", GetMetricHandler(memStorage))
	router.Get("/", GetPageHandler(memStorage))

	return router
}
