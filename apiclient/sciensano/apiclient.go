package sciensano

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/clambin/httpclient"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher"
	"github.com/go-http-utils/headers"
	"github.com/mailru/easyjson"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sort"
	"time"
)

const (
	TypeTestResults = iota
	TypeVaccinations
	TypeCases
	TypeMortality
	TypeHospitalisations
)

// Client calls the different Sciensano APIs
type Client struct {
	httpclient.Caller
	URL string
}

var _ fetcher.Fetcher = &Client{}

func (c Client) Fetch(ctx context.Context, dataType int) ([]apiclient.APIResponse, error) {
	resp, err := c.call(ctx, dataType, http.MethodGet)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return c.parseResponse(dataType, resp.Body)
}

func (c Client) GetLastUpdated(ctx context.Context, dataType int) (time.Time, error) {
	var lastModified time.Time
	resp, err := c.call(ctx, dataType, http.MethodHead)
	if err == nil {
		_ = resp.Body.Close()
		if lastModified, err = time.Parse(time.RFC1123, resp.Header.Get(headers.LastModified)); err != nil {
			err = fmt.Errorf("invalid timestamp %s: %w", headers.LastModified, err)
		}
	}
	return lastModified, err
}

func (c Client) DataTypes() map[int]string {
	return map[int]string{
		TypeCases:            "Cases",
		TypeHospitalisations: "Hospitalisations",
		TypeMortality:        "Mortality",
		TypeTestResults:      "TestResults",
		TypeVaccinations:     "Vaccinations",
	}
}

var endpoints = map[int]string{
	TypeCases:            "COVID19BE_CASES_AGESEX.json",
	TypeHospitalisations: "COVID19BE_HOSP.json",
	TypeMortality:        "COVID19BE_MORT.json",
	TypeTestResults:      "COVID19BE_tests.json",
	TypeVaccinations:     "COVID19BE_VACC.json",
}

func (c Client) getURL(dataType int) (string, error) {
	endpoint, found := endpoints[dataType]
	if !found {
		return "", fmt.Errorf("invalid data type: %d", dataType)
	}
	target := "https://epistat.sciensano.be"
	if c.URL != "" {
		target = c.URL
	}
	return target + "/Data/" + endpoint, nil
}

func (c Client) call(ctx context.Context, dataType int, method string) (resp *http.Response, err error) {
	var target string
	if target, err = c.getURL(dataType); err != nil {
		return
	}

	req, _ := http.NewRequestWithContext(ctx, method, target, nil)

	if resp, err = c.Caller.Do(req); err == nil && resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		err = fmt.Errorf("call failed: %s", resp.Status)
	}
	return
}

func (c Client) parseResponse(dataType int, body io.Reader) (results []apiclient.APIResponse, err error) {
	switch dataType {
	case TypeHospitalisations:
		var entries []*APIHospitalisationsResponse
		results, err = jsonDecode(body, entries)
	case TypeMortality:
		var entries []*APIMortalityResponse
		results, err = jsonDecode(body, entries)
	case TypeTestResults:
		var entries []*APITestResultsResponse
		results, err = jsonDecode(body, entries)
	case TypeCases:
		var entries APICasesResponses
		if err = easyjson.UnmarshalFromReader(body, &entries); err == nil {
			results = copyMaybeSort(entries)
		}
	case TypeVaccinations:
		var entries APIVaccinationsResponses
		if err = easyjson.UnmarshalFromReader(body, &entries); err == nil {
			results = copyMaybeSort(entries)
		}
	}
	return
}

func jsonDecode[T apiclient.APIResponse](body io.Reader, entries []T) (results []apiclient.APIResponse, err error) {
	if err = json.NewDecoder(body).Decode(&entries); err == nil {
		results = copyMaybeSort(entries)
	}
	return
}

// copyMaybeSort accomplishes two goals: it recasts the individual API response type to the generic APIResponse.
// Secondly, it checks if the data is in order. If not, it sorts it.
func copyMaybeSort[T apiclient.APIResponse](input []T) []apiclient.APIResponse {
	output := make([]apiclient.APIResponse, len(input))
	var timestamp time.Time
	var mustSort bool
	for idx, entry := range input {
		if !mustSort && entry.GetTimestamp().Before(timestamp) {
			mustSort = true
		}
		output[idx] = entry
		timestamp = entry.GetTimestamp()
	}
	if mustSort {
		log.Debug("sorting")
		sort.Slice(output, func(i, j int) bool { return output[i].GetTimestamp().Before(output[j].GetTimestamp()) })
		log.Debug("sorted")
	}
	return output
}
