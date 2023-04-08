package main

import (
	"context"
	"errors"
	"flag"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/simplejsonserver"
	"github.com/clambin/sciensano/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

var (
	debug            = flag.Bool("debug", false, "Log debug messages")
	simpleJSONAddr   = flag.String("addr", ":8080", "Server address")
	prometheusAddr   = flag.String("prometheus", ":9090", "Prometheus metrics port")
	demographicsPath = flag.String("demographics", "/data/population/TF_SOC_POP_STRUCT_2021.txt", "Path of the demographics server")
)

func main() {
	flag.Parse()

	var ops slog.HandlerOptions
	if *debug {
		ops.Level = slog.LevelDebug
		ops.AddSource = true
	}
	slog.SetDefault(slog.New(ops.NewTextHandler(os.Stdout)))

	slog.Info("Reporter API starting", "version", version.BuildVersion)

	go runPrometheusServer(*prometheusAddr)

	runSimpleJSONServer(*simpleJSONAddr, *demographicsPath)
}

func runPrometheusServer(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(addr, nil); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start Prometheus listener", "err", err)
	}
}

func runSimpleJSONServer(addr string, demographicsPath string) {
	s, err := simplejsonserver.New(&demographics.Server{
		Path:     demographicsPath,
		Interval: 24 * time.Hour,
	})
	if err != nil {
		slog.Error("failed to start SimpleJSON server", "err", err)
		panic(err)
	}
	prometheus.DefaultRegisterer.MustRegister(s)

	if err = s.Serve(context.Background(), addr); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start SimpleJSON server", "err", err)
		panic(err)
	}
}
