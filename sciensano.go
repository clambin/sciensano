package main

import (
	"github.com/clambin/sciensano/internal/apihandler"
	"github.com/clambin/sciensano/internal/version"
	"github.com/clambin/sciensano/pkg/grafana/apiserver"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.WithField("version", version.BuildVersion).Info("sciensano API starting")
	handler, _ := apihandler.Create()
	server := apiserver.Create(handler, 8080)
	_ = server.Run()
}
