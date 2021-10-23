package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"time"
)

// CasesGetter contains all methods providing COVID-19 cases
type CasesGetter interface {
	GetCases(ctx context.Context, endTime time.Time) (results []CaseCount, err error)
	GetCasesByProvince(ctx context.Context, endTime time.Time) (results map[string][]CaseCount, err error)
	GetCasesByRegion(ctx context.Context, endTime time.Time) (results map[string][]CaseCount, err error)
}

// CaseCount records the number of cases on a specific timestamp
type CaseCount struct {
	Timestamp time.Time
	Count     int
}

// GetCases returns all cases up to endTime
func (client *Client) GetCases(ctx context.Context, endTime time.Time) (results []CaseCount, err error) {
	var apiResult []*apiclient.APICasesResponse
	if apiResult, err = client.Getter.GetCases(ctx); err == nil {
		results = groupCases(apiResult, endTime)
	}
	return
}

func groupCases(cases []*apiclient.APICasesResponse, endTime time.Time) (totals []CaseCount) {
	var current CaseCount
	for _, entry := range cases {
		if entry.TimeStamp.IsZero() {
			continue
		}
		if entry.TimeStamp.After(endTime) {
			continue
		}
		if entry.TimeStamp.Time != current.Timestamp {
			if !current.Timestamp.IsZero() {
				totals = append(totals, current)
			}
			current = CaseCount{Timestamp: entry.TimeStamp.Time}
		}
		current.Count += entry.Cases
	}
	if !current.Timestamp.IsZero() {
		totals = append(totals, current)
	}
	return
}

// GetCasesByProvince returns all cases, grouped by province, up to endTime
func (client *Client) GetCasesByProvince(ctx context.Context, endTime time.Time) (results map[string][]CaseCount, err error) {
	var apiResult []*apiclient.APICasesResponse

	apiResult, err = client.Getter.GetCases(ctx)
	if err != nil {
		return
	}

	grouped := make(map[string][]*apiclient.APICasesResponse)
	for _, entry := range apiResult {
		grouped[entry.Province] = append(grouped[entry.Province], entry)
	}

	results = make(map[string][]CaseCount)
	for group := range grouped {
		results[group] = groupCases(grouped[group], endTime)
	}

	return
}

// GetCasesByRegion returns all cases, grouped by region, up to endTime
func (client *Client) GetCasesByRegion(ctx context.Context, endTime time.Time) (results map[string][]CaseCount, err error) {
	var apiResult []*apiclient.APICasesResponse

	apiResult, err = client.Getter.GetCases(ctx)
	if err != nil {
		return
	}

	grouped := make(map[string][]*apiclient.APICasesResponse)
	for _, entry := range apiResult {
		grouped[entry.Region] = append(grouped[entry.Region], entry)
	}

	results = make(map[string][]CaseCount)
	for group := range grouped {
		results[group] = groupCases(grouped[group], endTime)
	}

	return
}
