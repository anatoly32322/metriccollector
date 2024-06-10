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
		storage     *st.MemStorage
	}

	tests := []struct {
		name    string
		request string
		storage *st.MemStorage
		want    want
	}{
		{
			name: "simple test #1",
			want: want{
				contentType: "text/plain",
				statusCode:  200,
				storage: &st.MemStorage{
					GaugeMetrics:   map[string]float64{"Alloc": 1.2},
					CounterMetrics: make(map[string]int64),
					AcceptedMetricType: map[string]bool{
						"gauge":   true,
						"counter": true,
					},
				},
			},
			request: "/update/gauge/Alloc/1.2",
			storage: &st.MemStorage{
				GaugeMetrics:   make(map[string]float64),
				CounterMetrics: make(map[string]int64),
				AcceptedMetricType: map[string]bool{
					"gauge":   true,
					"counter": true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(MetricRouter(tt.storage, false, ""))
			defer ts.Close()

			request := httptest.NewRequest(http.MethodPost, tt.request, nil)

			resp, err := ts.Client().Do(request)
			if err != nil {
				return
			}

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.EqualValues(t, tt.want.storage, tt.storage)
			_ = resp.Body.Close()
		})
	}
}
