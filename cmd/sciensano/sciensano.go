package main

import (
	"context"
	"flag"
	"github.com/clambin/go-common/taskmanager"
	"github.com/clambin/go-common/taskmanager/httpserver"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/simplejsonserver"
	"github.com/clambin/sciensano/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
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

	var opts slog.HandlerOptions
	if *debug {
		opts.Level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &opts)))

	slog.Info("Reporter API starting", "version", version.BuildVersion)

	promHandler := http.NewServeMux()
	promHandler.Handle("/metrics", promhttp.Handler())

	tm := taskmanager.New(
		httpserver.New(*prometheusAddr, promHandler),
		getSimpleJSONServer(*simpleJSONAddr, *demographicsPath),
	)

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	err := tm.Run(ctx)
	if err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}
}

func getSimpleJSONServer(addr string, demographicsPath string) *httpserver.HTTPServer {
	s := simplejsonserver.New(&demographics.Server{
		Path:     demographicsPath,
		Interval: 24 * time.Hour,
	})
	prometheus.DefaultRegisterer.MustRegister(s)

	return httpserver.New(addr, s.Server)
}
