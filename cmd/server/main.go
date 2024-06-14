package main

import (
	"flag"
	"github.com/anatoly32322/metriccollector/internal/handlers"
	log "github.com/anatoly32322/metriccollector/internal/logger"
	st "github.com/anatoly32322/metriccollector/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Host            string `env:"ADDRESS"`
	StoreInterval   int64  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func gzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}

func main() {
	var err error
	var cfg Config
	flag.StringVar(&cfg.Host, "a", "localhost:8080", "hostname to listen on")
	flag.Int64Var(&cfg.StoreInterval, "i", 10, "interval to store metrics")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "path to file storage path")
	flag.BoolVar(&cfg.Restore, "r", true, "restore metrics from storage")

	flag.Parse()

	if envHostAddr := os.Getenv("ADDRESS"); envHostAddr != "" {
		cfg.Host = envHostAddr
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		cfg.StoreInterval, err = strconv.ParseInt(envStoreInterval, 10, 64)
		if err != nil {
			panic(err)
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		cfg.Restore, err = strconv.ParseBool(envRestore)
		if err != nil {
			panic(err)
		}
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	log.Sugar = *logger.Sugar()

	err = run(cfg)
	if err != nil {
		log.Sugar.Error(err)
	}
}

func run(cfg Config) error {
	memStorage := st.NewMemStorage()
	if cfg.Restore {
		err := memStorage.Load(cfg.FileStoragePath)
		if err != nil {
			return err
		}
	}

	defer func(memStorage *st.MemStorage, fname string) {
		err := memStorage.Save(fname)
		if err != nil {
			log.Sugar.Error(err)
		}
	}(memStorage, cfg.FileStoragePath)

	var router chi.Router

	if cfg.StoreInterval != 0 {
		go func() {
			var err error
			for {
				time.Sleep(time.Duration(cfg.StoreInterval) * time.Second)
				err = memStorage.Save(cfg.FileStoragePath)
				if err != nil {
					log.Sugar.Error(err)
				}
			}
		}()
		router = apihandlers.MetricRouter(memStorage, false, "")
	} else {
		router = apihandlers.MetricRouter(memStorage, true, cfg.FileStoragePath)
	}

	return http.ListenAndServe(cfg.Host, log.WithLogging(gzipMiddleware(router)))
}
