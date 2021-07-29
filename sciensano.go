package main

import (
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler"
	"github.com/clambin/sciensano/version"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.WithField("version", version.BuildVersion).Info("sciensano API starting")
	handler, _ := apihandler.Create()
	server := grafana_json.Create(handler, 8080)
	_ = server.Run()
}
