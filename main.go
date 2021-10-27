package main

import (
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"shinobi_exporter/config"
)

type Camera struct {
	Name string
}

var (
	argparser *flags.Parser
	opts      config.Opts
)

func initArgparser() {
	argparser = flags.NewParser(&opts, flags.Default)
	_, err := argparser.Parse()

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}

func main() {
	initArgparser()
	httpReqs := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "camera_hit_total",
			Help: "How many motion requests processed",
		},
		[]string{"cameraName"},
	)
	prometheus.MustRegister(httpReqs)
	var c Camera
	hitFunc := func(w http.ResponseWriter, r *http.Request) {
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
	log.Infof("Starting shinobi_exporter on %s", opts.ServerBind)
	if err := http.ListenAndServe(opts.ServerBind, nil); err != nil {
		log.WithField("error", err).Fatal("listen and serve error")
	}

}
