package apiserver

import (
	"encoding/json"
	"time"
)

// QueryRequest is the request passed to the Query API. It contains the list of targets to query, with
// timing & size parameters
type QueryRequest struct {
	From          time.Time
	To            time.Time
	MaxDataPoints uint64
	Targets       []string
}

// QueryResponse is the response from the Query API. It contains the target with its data points.
type QueryResponse struct {
	Target string
	Data   []QueryResponseData
}

type QueryResponseData struct {
	Timestamp time.Time
	Value     int64
}

func SaveSearchResponse(response []string) (body []byte, err error) {
	return json.Marshal(response)
}

func LoadQueryRequest(body []byte) (request QueryRequest, err error) {
	var req APIQueryRequest
	if err = json.Unmarshal(body, &req); err == nil {
		request = QueryRequest{
			From:          req.Range.From,
			To:            req.Range.To,
			MaxDataPoints: req.MaxDataPoints,
		}
		request.Targets = make([]string, len(req.Targets))
		for index, target := range req.Targets {
			request.Targets[index] = target.Target
		}
	}
	return
}

func SaveQueryResponse(responses []QueryResponse) (body []byte, err error) {
	resp := make([]APIQueryResponse, len(responses))
	for index, response := range responses {
		resp[index] = APIQueryResponse{
			Target: response.Target,
		}
		resp[index].DataPoints = make([][2]int64, len(response.Data))
		for index2, entry := range response.Data {
			resp[index].DataPoints[index2] = [2]int64{
				entry.Value,
				entry.Timestamp.UnixNano() / 1000000,
			}
		}
	}
	return json.Marshal(resp)
}
