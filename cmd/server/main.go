package main

import (
	"flag"
	"github.com/anatoly32322/metriccollector/internal/handlers"
	log "github.com/anatoly32322/metriccollector/internal/logger"
	st "github.com/anatoly32322/metriccollector/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func main() {
	host := flag.String("a", "localhost:8080", "hostname to listen on")
	flag.Parse()

	if envHostAddr := os.Getenv("ADDRESS"); envHostAddr != "" {
		host = &envHostAddr
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	log.Sugar = *logger.Sugar()

	run(*host)
}

func run(host string) {
	memStorage := st.NewMemStorage()

	router := apihandlers.MetricRouter(memStorage)

	err := http.ListenAndServe(host, log.WithLogging(router))
	if err != nil {
		return
	}
}
