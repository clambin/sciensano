package server

import (
	"encoding/json"
	"net/http"
	"slices"
)

func (s *Server) Health(w http.ResponseWriter, _ *http.Request) {
	dataSources := make([]string, 0, len(s.Handlers))
	for key := range s.Handlers {
		dataSources = append(dataSources, key)
	}
	slices.Sort(dataSources)
	response := struct {
		DataSources   []string
		ReporterCache []string
	}{
		DataSources:   dataSources,
		ReporterCache: s.reports.Keys(),
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(response)
}
