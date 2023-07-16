package main

import (
	"context"
	"flag"
	"github.com/clambin/go-common/taskmanager"
	"github.com/clambin/go-common/taskmanager/httpserver"
	promserver "github.com/clambin/go-common/taskmanager/prometheus"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/server"
	"github.com/clambin/sciensano/version"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
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

	s := server.New(&demographics.Server{
		Path:     *demographicsPath,
		Interval: 24 * time.Hour,
	})
	prometheus.DefaultRegisterer.MustRegister(s)

	tm := taskmanager.New(
		promserver.New(promserver.WithAddr(*prometheusAddr)),
		s,
		httpserver.New(*simpleJSONAddr, s.JSONServer),
	)

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	err := tm.Run(ctx)
	if err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}
}
