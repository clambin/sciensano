package fake

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

// Server runs a demographics test fake
type Server struct {
	Fail       bool
	filename   string
	httpServer *httptest.Server
}

// New creates a new test fake.  If filename is left blank, the standard test file is used
func New(filename string) (server *Server) {
	if filename == "" {
		filename = "../data/demographics.zip"
	}

	server = &Server{
		filename: filename,
	}
	server.httpServer = httptest.NewServer(http.HandlerFunc(server.handler))

	return
}

// Close closes the underlying httptest fake
func (server *Server) Close() {
	server.httpServer.Close()
}

// URL returns the URL of the underling httptest fake
func (server *Server) URL() string {
	return server.httpServer.URL
}

func (server *Server) handler(w http.ResponseWriter, _ *http.Request) {
	f, err := os.Open(server.filename)

	if err != nil {
		log.WithError(err).Error("handler failed")
		http.Error(w, "unable to open file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, f)

	if err != nil {
		log.WithError(err).Error("handler failed")
		http.Error(w, "unable to write data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_ = f.Close()
}
