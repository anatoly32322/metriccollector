package apihandlers

import (
	st "github.com/anatoly32322/metriccollector/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		storage     st.MemStorage
	}

	tests := []struct {
		name    string
		request string
		storage st.MemStorage
		want    want
	}{
		{
			name: "simple test #1",
			want: want{
				contentType: "text/plain",
				statusCode:  200,
				storage: st.MemStorage{
					GaugeMetrics:   st.GaugeMetric{Name: "Alloc", Value: 1.2},
					CounterMetrics: make([]st.CounterMetric, 0),
					AcceptedMetricType: map[string]bool{
						"gauge":   true,
						"counter": true,
					},
				},
			},
			request: "/update/gauge/Alloc/1.2",
			storage: st.MemStorage{
				GaugeMetrics:   st.GaugeMetric{},
				CounterMetrics: make([]st.CounterMetric, 0),
				AcceptedMetricType: map[string]bool{
					"gauge":   true,
					"counter": true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := ServeUpdateHandler(&tt.storage)
			request.Header.Add("Content-Type", "text/plain")

			h(w, request)

			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.EqualValues(t, tt.want.storage, tt.storage)
		})
	}
}
