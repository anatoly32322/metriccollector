package main

import (
	"github.com/anatoly32322/metriccollector/internal/handlers"
	st "github.com/anatoly32322/metriccollector/internal/storage"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	run()
}

func run() {
	memStorage := st.NewMemStorage()

	router := apihandlers.MetricRouter(memStorage)
	log.Fatal(http.ListenAndServe(":8080", router))
}
