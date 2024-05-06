package main

import (
	"github.com/anatoly32322/metriccollector/internal/handlers"
	st "github.com/anatoly32322/metriccollector/internal/storage"
	"net/http"
)

func main() {
	run()
}

func run() {
	memStorage := st.NewMemStorage()
	mux := http.NewServeMux()

	mux.Handle("/update/", apihandlers.ServeUpdateHandler(memStorage))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}
