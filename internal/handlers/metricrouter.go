package apihandlers

import (
	log "github.com/anatoly32322/metriccollector/internal/logger"
	st "github.com/anatoly32322/metriccollector/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func syncStoreMiddleware(memStorage st.Storage, storePath string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			err := memStorage.Save(storePath)
			if err != nil {
				log.Sugar.Error(err)
			}
		})
	}
}

func MetricRouter(memStorage st.Storage, isSyncStore bool, storePath string) chi.Router {
	router := chi.NewRouter()
	updateSubRouter := chi.NewRouter()

	router.Use(middleware.Compress(5))

	if isSyncStore {
		storeMiddleware := syncStoreMiddleware(memStorage, storePath)
		updateSubRouter.Use(storeMiddleware)
	}
	updateSubRouter.Post("/{metricType}/{metricName}/{metricValue}", ServeUpdateHandler(memStorage))
	updateSubRouter.Post("/", ServeUpdateHandlerV2(memStorage))

	router.Mount("/update", updateSubRouter)
	router.Post("/value/", GetMetricHandlerV2(memStorage))
	router.Get("/value/{metricType}/{metricName}", GetMetricHandler(memStorage))
	router.Get("/", GetPageHandler(memStorage))

	return router
}
