package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/simplejsonserver"
	"github.com/clambin/sciensano/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var (
		debug            bool
		simpleJSONPort   int
		prometheusPort   int
		demographicsPath string
	)

	a := kingpin.New(filepath.Base(os.Args[0]), "reporter")
	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").Short('d').BoolVar(&debug)
	a.Flag("port", "Server port").Short('p').Default("8080").IntVar(&simpleJSONPort)
	a.Flag("prometheus", "Prometheus metrics port").Default("9090").IntVar(&prometheusPort)
	a.Flag("demographics", "Path of the demographics server").Default("/data/population/TF_SOC_POP_STRUCT_2021.txt").StringVar(&demographicsPath)

	if _, err := a.Parse(os.Args[1:]); err != nil {
		a.Usage(os.Args[1:])
		os.Exit(1)
	}

	var ops slog.HandlerOptions
	if debug {
		ops.Level = slog.LevelDebug
		ops.AddSource = true
	}
	slog.SetDefault(slog.New(ops.NewTextHandler(os.Stdout)))

	slog.Info("Reporter API starting", "version", version.BuildVersion)

	go runPrometheusServer(prometheusPort)

	runSimpleJSONServer(simpleJSONPort, demographicsPath)
}

func runPrometheusServer(port int) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start Prometheus listener", "err", err)
	}
}

func runSimpleJSONServer(port int, demographicsPath string) {
	s, err := simplejsonserver.New(&demographics.Server{
		Path:     demographicsPath,
		Interval: 24 * time.Hour,
	})
	if err != nil {
		slog.Error("failed to start SimpleJSON server", "err", err)
		panic(err)
	}
	prometheus.DefaultRegisterer.MustRegister(s)

	if err = s.Serve(context.Background(), port); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start SimpleJSON server", "err", err)
		panic(err)
	}
}
