package main

import (
	"context"
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/version"
	log "github.com/sirupsen/logrus"
	"net/http"
	// _ "net/http/pprof"
	"time"
)

func main() {
	log.WithField("version", version.BuildVersion).Info("sciensano API starting")
	demo := demographics.New()
	go func() {
		_ = demo.Run(context.Background(), 24*time.Hour)
	}()
	server := grafana_json.Create(apihandler.Create(demo))
	r := server.GetRouter()
	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(r)
	}
}
