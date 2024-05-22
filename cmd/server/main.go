package main

import (
	"flag"
	"github.com/anatoly32322/metriccollector/internal/handlers"
	st "github.com/anatoly32322/metriccollector/internal/storage"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	host := flag.String("a", "localhost:8080", "hostname to listen on")
	flag.Parse()

	if envHostAddr := os.Getenv("ADDRESS"); envHostAddr != "" {
		host = &envHostAddr
	}

	run(*host)
}

func run(host string) {
	memStorage := st.NewMemStorage()

	router := apihandlers.MetricRouter(memStorage)
	log.Info(host)
	log.Fatal(http.ListenAndServe(host, router))
}
