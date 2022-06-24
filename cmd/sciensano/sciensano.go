package main

import (
	"errors"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver"
	"github.com/clambin/sciensano/version"
	"github.com/clambin/simplejson/v3"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var (
		debug            bool
		port             int
		demographicsPath string
	)

	log.WithField("version", version.BuildVersion).Info("Reporter API starting")
	a := kingpin.New(filepath.Base(os.Args[0]), "reporter")
	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").Short('d').BoolVar(&debug)
	a.Flag("port", "Server port").Short('p').Default("8080").IntVar(&port)
	a.Flag("demographics", "Path of the demographics server").Default("/data/population/TF_SOC_POP_STRUCT_2021.txt").StringVar(&demographicsPath)

	if _, err := a.Parse(os.Args[1:]); err != nil {
		a.Usage(os.Args[1:])
		os.Exit(1)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	s := simplejsonserver.Server{
		Server:   simplejson.Server{Name: "sciensano"},
		Reporter: reporter.New(15 * time.Minute),
		Demographics: &demographics.Server{
			Path:     demographicsPath,
			Interval: 24 * time.Hour,
		},
	}

	if err := s.Run(port); !errors.Is(err, http.ErrServerClosed) {
		log.WithError(err).Fatal("failed to start HTTP server")
	}
}
