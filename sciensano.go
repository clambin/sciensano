package main

import (
	"github.com/clambin/covid19/pkg/grafana/apiserver"
	"github.com/clambin/sciensano/internal/apihandler"
	"github.com/clambin/sciensano/internal/version"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.WithField("version", version.BuildVersion).Info("sciensano API starting")
	handler, _ := apihandler.Create()
	server := apiserver.Create(handler, 8080)
	_ = server.Run()
}
