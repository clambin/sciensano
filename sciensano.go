package main

import (
	"github.com/clambin/sciensano/apihandler"
	"github.com/clambin/sciensano/version"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
)

func main() {
	var (
		debug bool
		port  int
	)

	log.WithField("version", version.BuildVersion).Info("Sciensano API starting")
	a := kingpin.New(filepath.Base(os.Args[0]), "reporter")
	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").Short('d').BoolVar(&debug)
	a.Flag("port", "Server port").Short('p').Default("8080").IntVar(&port)

	if _, err := a.Parse(os.Args[1:]); err != nil {
		a.Usage(os.Args[1:])
		os.Exit(1)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	handler := apihandler.NewServer()

	if err := handler.Run(port); err != http.ErrServerClosed {
		log.WithError(err).Fatal("failed to start HTTP server")
	}
}
