package apiserver

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func (apiServer *APIServer) hello(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Hello")
}

func (apiServer *APIServer) search(w http.ResponseWriter, _ *http.Request) {
	response := apiServer.apiHandler.Search()
	if targetsJSON, err := SaveSearchResponse(response); err == nil {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(targetsJSON)
	} else {
		log.WithField("err", err).Warning("failed to create search response")
		http.Error(w, "failed to create search response", http.StatusInternalServerError)
	}
}

func (apiServer *APIServer) query(w http.ResponseWriter, req *http.Request) {
	var (
		err      error
		bytes    []byte
		request  QueryRequest
		response []QueryResponse
	)
	defer req.Body.Close()

	if bytes, err = ioutil.ReadAll(req.Body); err == nil {
		request, err = LoadQueryRequest(bytes)
	}

	if err == nil {
		if response, err = apiServer.apiHandler.Query(&request); err == nil {
			bytes, err = SaveQueryResponse(response)
		}

		if err == nil {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(bytes)
		} else {
			log.WithField("err", err).Warning("failed to process query request")
			http.Error(w, "failed to process query request", http.StatusInternalServerError)
		}
	} else {
		log.WithField("err", err).Warning("failed to parse query request")
		http.Error(w, "failed to parse query request", http.StatusBadRequest)
	}
}
