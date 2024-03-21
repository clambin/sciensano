package main

import (
	"context"
	"errors"
	"flag"
	"github.com/clambin/go-common/http/roundtripper"
	"github.com/clambin/go-common/taskmanager"
	"github.com/clambin/go-common/taskmanager/httpserver"
	promserver "github.com/clambin/go-common/taskmanager/prometheus"
	gjson "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/v2/internal/population"
	"github.com/clambin/sciensano/v2/internal/reports"
	"github.com/clambin/sciensano/v2/internal/reports/datasource"
	"github.com/clambin/sciensano/v2/internal/reports/store"
	"github.com/clambin/sciensano/v2/internal/server"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"
)

var (
	version = "change-me"

	debug            = flag.Bool("debug", false, "Log debug messages")
	simpleJSONAddr   = flag.String("addr", ":8080", "Server address")
	prometheusAddr   = flag.String("prometheus", ":9090", "Prometheus metrics port")
	demographicsPath = flag.String("demographics", "/data/population/TF_SOC_POP_STRUCT_2023.txt", "Path of the demographics file")
)

func main() {
	flag.Parse()

	var opts slog.HandlerOptions
	if *debug {
		opts.Level = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &opts))

	logger.Info("Sciensano API server starting", "version", version)

	popStore := population.Server{Path: *demographicsPath, Interval: 24 * time.Hour, Logger: logger.With("component", "population")}

	reportsStore := store.Store{Logger: logger.With("component", "reportsStore")}

	httpMetrics := roundtripper.NewDefaultRoundTripMetrics("sciensano", "", "sciensano")
	prometheus.MustRegister(httpMetrics)
	r := roundtripper.New(
		roundtripper.WithLimiter(3),
		roundtripper.WithInstrumentedRoundTripper(httpMetrics),
	)
	client := &http.Client{Transport: r}

	ds := datasource.NewSciensanoDatastore("", 15*time.Minute, client, logger.With("component", "datasource"))
	reporters := reports.NewSciensanoReporters(ds, &reportsStore, &popStore, logger.With("component", "reporters"))

	var tasks []taskmanager.Task
	tasks = append(tasks, ds)
	tasks = append(tasks, &popStore)
	tasks = append(tasks, reporters...)

	gjsonMetrics := gjson.NewDefaultPrometheusQueryMetrics("sciensano", "", "sciensano")
	prometheus.MustRegister(gjsonMetrics)
	s := server.New(&reportsStore, gjsonMetrics, logger.With("component", "server"))

	tasks = append(
		tasks, promserver.New(promserver.WithAddr(*prometheusAddr)),
		s,
		httpserver.New(*simpleJSONAddr, s.JSONServer),
		httpserver.New(":6060", http.DefaultServeMux),
	)
	tm := taskmanager.New(tasks...)

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	if err := tm.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("failed to start", "err", err)
		os.Exit(1)
	}
}
