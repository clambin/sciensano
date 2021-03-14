package apiserver

import "time"

// APIQueryRequest contains the request parameters to the API's 'query' method
type APIQueryRequest struct {
	Range         APIQueryRequestRange    `json:"range"`
	Interval      string                  `json:"interval"`
	MaxDataPoints uint64                  `json:"maxDataPoints"`
	Targets       []APIQueryRequestTarget `json:"targets"`
}

type APIQueryRequestRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type APIQueryRequestTarget struct {
	Target string `json:"target"`
	Type   string `json:"type"`
}

// APIQueryResponse contains the response of the API's 'query' method
type APIQueryResponse struct {
	Target     string     `json:"target"`
	DataPoints [][2]int64 `json:"datapoints"`
}
