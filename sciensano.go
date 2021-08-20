package main

import (
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler"
	"github.com/clambin/sciensano/version"
	log "github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	log.WithField("version", version.BuildVersion).Info("sciensano API starting")
	server := grafana_json.Create(apihandler.Create())
	r := server.GetRouter()
	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(r)
	}
}
