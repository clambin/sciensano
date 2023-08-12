package main

import (
	"context"
	"flag"
	"github.com/clambin/go-common/taskmanager"
	"github.com/clambin/go-common/taskmanager/httpserver"
	promserver "github.com/clambin/go-common/taskmanager/prometheus"
	"github.com/clambin/sciensano/internal/population"
	"github.com/clambin/sciensano/internal/reports/datasource"
	"github.com/clambin/sciensano/internal/reports/reporter"
	"github.com/clambin/sciensano/internal/reports/store"
	"github.com/clambin/sciensano/internal/server"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"
)

var (
	// BuildVersion contains the release number
	BuildVersion = "change-me"

	debug            = flag.Bool("debug", false, "Log debug messages")
	simpleJSONAddr   = flag.String("addr", ":8080", "Server address")
	prometheusAddr   = flag.String("prometheus", ":9090", "Prometheus metrics port")
	demographicsPath = flag.String("demographics", "/data/population/TF_SOC_POP_STRUCT_2021.txt", "Path of the demographics server")
)

func main() {
	flag.Parse()

	var opts slog.HandlerOptions
	if *debug {
		opts.Level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &opts)))

	slog.Info("Sciensano API server starting", "version", BuildVersion)

	popStore := population.Server{Path: *demographicsPath, Interval: 24 * time.Hour}

	reportsStore := store.Store{Logger: slog.Default().With("component", "reportsStore")}
	ds := datasource.NewSciensanoDatastore("", 15*time.Minute, slog.Default().With("component", "datasource"))
	reporters := reporter.NewSciensanoReporters(ds, &reportsStore, &popStore, slog.Default().With("component", "reporters"))

	var tasks []taskmanager.Task
	tasks = append(tasks, ds)
	tasks = append(tasks, &popStore)
	tasks = append(tasks, reporters...)

	s := server.New(&reportsStore, slog.Default())
	prometheus.DefaultRegisterer.MustRegister(s)

	tasks = append(
		tasks, promserver.New(promserver.WithAddr(*prometheusAddr)),
		s,
		httpserver.New(*simpleJSONAddr, s.JSONServer),
		httpserver.New(":6060", http.DefaultServeMux),
	)
	tm := taskmanager.New(tasks...)

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	err := tm.Run(ctx)
	if err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}
}
