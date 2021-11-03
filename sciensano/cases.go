package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/datasets"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

// CasesGetter contains all methods providing COVID-19 cases
type CasesGetter interface {
	GetCases(ctx context.Context) (results *datasets.Dataset, err error)
	GetCasesByRegion(ctx context.Context) (results *datasets.Dataset, err error)
	GetCasesByProvince(ctx context.Context) (results *datasets.Dataset, err error)
	GetCasesByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error)
}

const (
	groupCasesByNone = iota
	groupCasesByRegion
	groupCasesByProvince
	groupCasesByAgeGroup
)

// GetCases returns all cases
func (client *Client) GetCases(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getCases(ctx, "GetCases", "Cases", groupCasesByNone)
}

// GetCasesByRegion returns all cases, grouped by region
func (client *Client) GetCasesByRegion(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getCases(ctx, "GetCasesByRegion", "CasesByRegion", groupCasesByRegion)
}

// GetCasesByProvince returns all cases, grouped by province
func (client *Client) GetCasesByProvince(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getCases(ctx, "GetCasesByProvince", "CasesByProvince", groupCasesByProvince)
}

// GetCasesByAgeGroup returns all cases, grouped by province
func (client *Client) GetCasesByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getCases(ctx, "GetCasesByAgeGroup", "CasesByAgeGroup", groupCasesByAgeGroup)
}

// getCases
func (client *Client) getCases(ctx context.Context, name, cacheEntryName string, mode int) (results *datasets.Dataset, err error) {
	before := time.Now()
	defer func() { log.WithField("time", time.Now().Sub(before)).Debug(name + " done") }()

	log.Debug("running " + name)
	entry := client.cache.Load(cacheEntryName)
	entry.Once.Do(func() {
		var apiResult apiclient.APICasesResponse
		if apiResult, err = client.Getter.GetCases(ctx); err == nil {
			entry.Data = groupCases(apiResult, mode)
			client.cache.Save(cacheEntryName, entry)
		} else {
			client.cache.Clear(cacheEntryName)
		}
	})
	if err == nil && entry.Data != nil {
		results = entry.Data.Copy()
	}
	return
}

func groupCases(cases apiclient.APICasesResponse, groupField int) (results *datasets.Dataset) {
	mappedCases, mappedGroups := mapCases(cases, groupField)
	timestamps := getUniqueSortedTimestampsFromCases(mappedCases)
	groups := getUniqueSortedGroupNames(mappedGroups)

	results = &datasets.Dataset{
		Timestamps: timestamps,
		Groups:     make([]datasets.GroupedDatasetEntry, len(groups)),
	}

	for index, group := range groups {
		results.Groups[index] = datasets.GroupedDatasetEntry{
			Name: group,
		}
	}

	for _, timestamp := range timestamps {
		for index, group := range groups {
			value, ok := mappedCases[timestamp][group]
			if ok == false {
				value = &CasesEntry{}
			}
			results.Groups[index].Values = append(results.Groups[index].Values, value)
		}
	}
	return

}

func mapCases(cases apiclient.APICasesResponse, groupField int) (mappedCases map[time.Time]map[string]*CasesEntry, mappedGroups map[string]struct{}) {
	mappedCases = make(map[time.Time]map[string]*CasesEntry)
	mappedGroups = make(map[string]struct{})

	for _, entry := range cases {
		mappedCase, ok := mappedCases[entry.TimeStamp.Time]
		if ok == false {
			mappedCase = make(map[string]*CasesEntry)
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

		value, ok := mappedCase[groupName]
		if ok == false {
			value = &CasesEntry{}
		}

		value.Count += entry.Cases
		mappedCase[groupName] = value
		mappedCases[entry.TimeStamp.Time] = mappedCase

		mappedGroups[groupName] = struct{}{}
	}

	return
}

func getUniqueSortedTimestampsFromCases(input map[time.Time]map[string]*CasesEntry) (output []time.Time) {
	for timestamp := range input {
		output = append(output, timestamp)
	}
	sort.Slice(output, func(i, j int) bool { return output[i].Before(output[j]) })
	return
}
