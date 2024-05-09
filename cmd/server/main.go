package main

import (
	"flag"
	"github.com/anatoly32322/metriccollector/internal/handlers"
	st "github.com/anatoly32322/metriccollector/internal/storage"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	host := flag.String("a", "localhost:8080", "hostname to listen on")
	flag.Parse()
	run(*host)
}

func run(host string) {
	memStorage := st.NewMemStorage()

	router := apihandlers.MetricRouter(memStorage)
	log.Info(host)
	log.Fatal(http.ListenAndServe(host, router))
}
