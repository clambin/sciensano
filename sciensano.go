package main

import (
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/apihandler"
	"github.com/clambin/sciensano/internal/version"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.WithField("version", version.BuildVersion).Info("sciensano API starting")
	handler, _ := apihandler.Create()
	server := grafana_json.Create(handler.Endpoints, 8080)
	_ = server.Run()
}
