package apihandlers

import (
	st "github.com/anatoly32322/metriccollector/internal/storage"
	"log"
	"net/http"
	"strings"
)

func ServeUpdateHandler(memStorage *st.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		//if r.FormValue("Content-Type") != "text/plain" {
		//	w.WriteHeader(http.StatusBadRequest)
		//	return
		//}

		pathParts := strings.Split(r.URL.Path, "/")

		log.Print(pathParts)

		if len(pathParts) != 5 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metricType := pathParts[2]
		metricName := pathParts[3]
		value := pathParts[4]
		err := memStorage.Update(metricType, metricName, value)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		return
	}
}
