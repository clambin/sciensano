package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"sort"
	"time"
)

// CasesGetter contains all methods providing COVID-19 cases
type CasesGetter interface {
	GetCases(ctx context.Context, endTime time.Time) (results *Cases, err error)
	GetCasesByRegion(ctx context.Context, endTime time.Time) (results *Cases, err error)
	GetCasesByProvince(ctx context.Context, endTime time.Time) (results *Cases, err error)
	GetCasesByAgeGroup(ctx context.Context, endTime time.Time) (results *Cases, err error)
}

// Cases is the response of the GetCases functions
type Cases struct {
	Timestamps []time.Time
	Groups     []GroupedCasesEntry
}

// GroupedCasesEntry contains all the values for the (grouped) cases
type GroupedCasesEntry struct {
	Name   string
	Values []int
}

const (
	groupCasesByNone = iota
	groupCasesByRegion
	groupCasesByProvince
	groupCasesByAgeGroup
)

// GetCases returns all cases up to endTime
func (client *Client) GetCases(ctx context.Context, endTime time.Time) (results *Cases, err error) {
	var apiResult []*apiclient.APICasesResponse
	if apiResult, err = client.Getter.GetCases(ctx); err != nil {
		return
	}

	return groupCases(apiResult, endTime, groupCasesByNone), nil
}

// GetCasesByRegion returns all cases up to endTime, grouped by region
func (client *Client) GetCasesByRegion(ctx context.Context, endTime time.Time) (results *Cases, err error) {
	var apiResult []*apiclient.APICasesResponse
	if apiResult, err = client.Getter.GetCases(ctx); err != nil {
		return
	}

	return groupCases(apiResult, endTime, groupCasesByRegion), nil
}

// GetCasesByProvince returns all cases up to endTime, grouped by province
func (client *Client) GetCasesByProvince(ctx context.Context, endTime time.Time) (results *Cases, err error) {
	var apiResult []*apiclient.APICasesResponse
	if apiResult, err = client.Getter.GetCases(ctx); err != nil {
		return
	}

	return groupCases(apiResult, endTime, groupCasesByProvince), nil
}

// GetCasesByAgeGroup returns all cases up to endTime, grouped by province
func (client *Client) GetCasesByAgeGroup(ctx context.Context, endTime time.Time) (results *Cases, err error) {
	var apiResult []*apiclient.APICasesResponse
	if apiResult, err = client.Getter.GetCases(ctx); err != nil {
		return
	}

	return groupCases(apiResult, endTime, groupCasesByAgeGroup), nil
}

func groupCases(cases []*apiclient.APICasesResponse, endTime time.Time, groupField int) (results *Cases) {
	mappedCases, mappedGroups := mapCases(cases, endTime, groupField)
	timestamps := getUniqueSortedTimestampsFromCases(mappedCases)
	groups := getUniqueSortedGroupNames(mappedGroups)

	results = &Cases{
		Timestamps: timestamps,
		Groups:     make([]GroupedCasesEntry, len(groups)),
	}

	for index, group := range groups {
		results.Groups[index] = GroupedCasesEntry{
			Name: group,
		}
	}

	for _, timestamp := range timestamps {
		for index, group := range groups {
			value, _ := mappedCases[timestamp][group]
			results.Groups[index].Values = append(results.Groups[index].Values, value)
		}
	}
	return

}

func mapCases(cases []*apiclient.APICasesResponse, endTime time.Time, groupField int) (mappedCases map[time.Time]map[string]int, mappedGroups map[string]struct{}) {
	mappedCases = make(map[time.Time]map[string]int)
	mappedGroups = make(map[string]struct{})

	for _, entry := range cases {
		if entry.TimeStamp.After(endTime) {
			continue
		}

		mappedCase, ok := mappedCases[entry.TimeStamp.Time]
		if ok == false {
			mappedCase = make(map[string]int)
		}

		var groupName string
		switch groupField {
		case groupCasesByNone:
			groupName = ""
		case groupCasesByRegion:
			groupName = entry.Region
		case groupCasesByProvince:
			groupName = entry.Province
		case groupCasesByAgeGroup:
			groupName = entry.AgeGroup
		}

		value, _ := mappedCase[groupName]
		value += entry.Cases
		mappedCase[groupName] = value
		mappedCases[entry.TimeStamp.Time] = mappedCase

		mappedGroups[groupName] = struct{}{}
	}

	return
}

func getUniqueSortedTimestampsFromCases(input map[time.Time]map[string]int) (output []time.Time) {
	for timestamp := range input {
		output = append(output, timestamp)
	}
	sort.Slice(output, func(i, j int) bool { return output[i].Before(output[j]) })
	return
}
