package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) Health(w http.ResponseWriter, _ *http.Request) {
	response := struct {
		DataSources   int
		ReporterCache []string
	}{
		DataSources:   len(s.Handlers),
		ReporterCache: s.reports.Keys(),
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(response)
}
