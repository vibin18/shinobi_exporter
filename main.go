package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)
type Camera struct {
	Name string
}

func main() {
	httpReqs := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "camera_hit_total",
			Help: "How many motion requests processed",
		},
		[]string{"cameraName"},
	)
	prometheus.MustRegister(httpReqs)
	var c Camera
	hitFunc := func (w http.ResponseWriter, r *http.Request){
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.WithField("error", err).Error("error handling request")
			return
		}
		httpReqs.WithLabelValues(c.Name).Inc()
	}

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hit", hitFunc)
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.WithField("error", err).Fatal("listen and serve error")
	}

}