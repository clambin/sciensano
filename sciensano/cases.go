package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"sort"
	"time"
)

// CasesGetter contains all methods providing COVID-19 cases
type CasesGetter interface {
	GetCases(ctx context.Context, endTime time.Time) (results *GroupedCases, err error)
	GetCasesByRegion(ctx context.Context, endTime time.Time) (results *GroupedCases, err error)
	GetCasesByProvince(ctx context.Context, endTime time.Time) (results *GroupedCases, err error)
	GetCasesByAgeGroup(ctx context.Context, endTime time.Time) (results *GroupedCases, err error)
}

// GroupedCases is the response of the GetCases functions
type GroupedCases struct {
	Timestamps []time.Time
	Groups     []GroupedCasesEntry
}

// GroupedCasesEntry contains all the values for the (grouped) cases
type GroupedCasesEntry struct {
	Name   string
	Values []int
}

const (
	groupByNone = iota
	groupByRegion
	groupByProvince
	groupByAgeGroup
)

// GetCases returns all cases up to endTime
func (client *Client) GetCases(ctx context.Context, endTime time.Time) (results *GroupedCases, err error) {
	var apiResult []*apiclient.APICasesResponse
	if apiResult, err = client.Getter.GetCases(ctx); err != nil {
		return
	}

	return groupCases(apiResult, endTime, groupByNone), nil
}

// GetCasesByRegion returns all cases up to endTime, grouped by region
func (client *Client) GetCasesByRegion(ctx context.Context, endTime time.Time) (results *GroupedCases, err error) {
	var apiResult []*apiclient.APICasesResponse
	if apiResult, err = client.Getter.GetCases(ctx); err != nil {
		return
	}

	return groupCases(apiResult, endTime, groupByRegion), nil
}

// GetCasesByProvince returns all cases up to endTime, grouped by province
func (client *Client) GetCasesByProvince(ctx context.Context, endTime time.Time) (results *GroupedCases, err error) {
	var apiResult []*apiclient.APICasesResponse
	if apiResult, err = client.Getter.GetCases(ctx); err != nil {
		return
	}

	return groupCases(apiResult, endTime, groupByProvince), nil
}

// GetCasesByAgeGroup returns all cases up to endTime, grouped by province
func (client *Client) GetCasesByAgeGroup(ctx context.Context, endTime time.Time) (results *GroupedCases, err error) {
	var apiResult []*apiclient.APICasesResponse
	if apiResult, err = client.Getter.GetCases(ctx); err != nil {
		return
	}

	return groupCases(apiResult, endTime, groupByAgeGroup), nil
}

func groupCases(cases []*apiclient.APICasesResponse, endTime time.Time, indexField int) (results *GroupedCases) {
	mappedCases := make(map[time.Time]map[string]int)
	mappedGroups := make(map[string]struct{})

	for _, entry := range cases {
		if entry.TimeStamp.After(endTime) {
			continue
		}

		mappedCase, ok := mappedCases[entry.TimeStamp.Time]
		if ok == false {
			mappedCase = make(map[string]int)
		}

		var groupName string
		switch indexField {
		case groupByNone:
			groupName = ""
		case groupByRegion:
			groupName = entry.Region
		case groupByProvince:
			groupName = entry.Province
		case groupByAgeGroup:
			groupName = entry.AgeGroup
		}

		value, _ := mappedCase[groupName]
		value += entry.Cases
		mappedCase[groupName] = value
		mappedCases[entry.TimeStamp.Time] = mappedCase

		mappedGroups[groupName] = struct{}{}
	}

	timestamps := getUniqueSortedTimestamps(mappedCases)
	groups := getUniqueSortedGroupNames(mappedGroups)

	results = &GroupedCases{
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

func getUniqueSortedTimestamps(input map[time.Time]map[string]int) (output []time.Time) {
	for timestamp := range input {
		output = append(output, timestamp)
	}
	sort.Slice(output, func(i, j int) bool { return output[i].Before(output[j]) })
	return
}

func getUniqueSortedGroupNames(input map[string]struct{}) (output []string) {
	for group := range input {
		output = append(output, group)
	}
	sort.Strings(output)
	return
}
