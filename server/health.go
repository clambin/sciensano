package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) Health(w http.ResponseWriter, _ *http.Request) {
	response := struct {
		DataSources   int
		ReporterCache map[string]int
	}{
		DataSources: len(s.Handlers),
		//ReporterCache: s.reportsCache.Stats(),
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(response)
}
