package main

import (
	"context"
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/version"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	log.WithField("version", version.BuildVersion).Info("sciensano API starting")
	demo := demographics.New()
	go func() {
		_ = demo.Run(context.Background(), 24*time.Hour)
	}()
	handler, _ := apihandler.Create(demo)
	server := grafana_json.Create(handler, 8080)
	_ = server.Run()
}
