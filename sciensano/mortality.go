package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/datasets"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

// MortalityGetter contains all methods providing COVID-19 mortality
type MortalityGetter interface {
	GetMortality(ctx context.Context) (results *datasets.Dataset, err error)
	GetMortalityByRegion(ctx context.Context) (results *datasets.Dataset, err error)
	GetMortalityByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error)
}

// GroupedMortalityEntry contains all the values for the (grouped) mortality figures
type GroupedMortalityEntry struct {
	Name   string
	Values []int
}

const (
	groupMortalityByNone = iota
	groupMortalityByRegion
	groupMortalityByAgeGroup
)

// GetMortality returns all mortality figures
func (client *Client) GetMortality(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getMortality(ctx, "GetMortality", "Mortality", groupMortalityByNone)
}

// GetMortalityByRegion returns all mortality figures, grouped by region
func (client *Client) GetMortalityByRegion(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getMortality(ctx, "GetMortalityByRegion", "MortalityByRegion", groupMortalityByRegion)
}

// GetMortalityByAgeGroup returns all Mortality, grouped by age group
func (client *Client) GetMortalityByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getMortality(ctx, "GetMortalityByAge", "MortalityByAge", groupMortalityByAgeGroup)
}

func (client *Client) getMortality(ctx context.Context, name, cacheEntryName string, mode int) (results *datasets.Dataset, err error) {
	before := time.Now()
	defer func() { log.WithField("time", time.Now().Sub(before)).Debug(name + " done") }()

	log.Debug("running " + name)
	entry := client.cache.Load(cacheEntryName)
	entry.Once.Do(func() {
		var apiResult apiclient.APIMortalityResponse
		if apiResult, err = client.Getter.GetMortality(ctx); err == nil {
			entry.Data = groupMortality(apiResult, mode)
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

func groupMortality(mortality apiclient.APIMortalityResponse, groupField int) (results *datasets.Dataset) {
	before := time.Now()
	defer log.WithField("time", time.Now().Sub(before)).Debug("groupMortality")

	mappedMortalities, mappedGroups := mapMortality(mortality, groupField)
	timestamps := getUniqueSortedTimestampsFromMortality(mappedMortalities)
	groups := getUniqueSortedGroupNames(mappedGroups)

	results = &datasets.Dataset{
		Timestamps: timestamps,
		Groups:     make([]datasets.GroupedDatasetEntry, len(groups)),
	}

	for index, group := range groups {
		results.Groups[index] = datasets.GroupedDatasetEntry{
			Name:   group,
			Values: make([]datasets.Copyable, 0, len(timestamps)),
		}
	}

	for _, timestamp := range timestamps {
		for index, group := range groups {
			entry, ok := mappedMortalities[timestamp][group]
			if ok == false {
				entry = &MortalityEntry{}
			}
			results.Groups[index].Values = append(results.Groups[index].Values, entry)
		}
	}
	return
}

func mapMortality(mortality apiclient.APIMortalityResponse, groupField int) (mappedMortality map[time.Time]map[string]*MortalityEntry, mappedGroups map[string]struct{}) {
	before := time.Now()
	defer log.WithField("time", time.Now().Sub(before)).Debug("mapMortality")

	if len(mortality) == 0 {
		return
	}

	mappedMortality = make(map[time.Time]map[string]*MortalityEntry)
	mappedGroups = make(map[string]struct{})

	currentTimestamp := time.Time{}
	currentEntries := make(map[string]*MortalityEntry)

	for _, mortalityFigure := range mortality {
		if !currentTimestamp.IsZero() && !currentTimestamp.Equal(mortalityFigure.TimeStamp.Time) {
			mappedMortality[currentTimestamp] = currentEntries
			currentEntries = make(map[string]*MortalityEntry)
		}
		currentTimestamp = mortalityFigure.TimeStamp.Time

		var groupName string
		switch groupField {
		case groupMortalityByNone:
			groupName = ""
		case groupMortalityByAgeGroup:
			groupName = mortalityFigure.AgeGroup
		case groupMortalityByRegion:
			groupName = mortalityFigure.Region
		}

		mappedGroups[groupName] = struct{}{}

		entry, ok := currentEntries[groupName]

		if ok == false {
			entry = &MortalityEntry{}
		}

		entry.Count += mortalityFigure.Deaths
		currentEntries[groupName] = entry
	}

	if !currentTimestamp.IsZero() {
		mappedMortality[currentTimestamp] = currentEntries
	}
	return
}

func getUniqueSortedTimestampsFromMortality(input map[time.Time]map[string]*MortalityEntry) (output []time.Time) {
	for timestamp := range input {
		output = append(output, timestamp)
	}
	sort.Slice(output, func(i, j int) bool { return output[i].Before(output[j]) })
	return
}

type MortalityEntry struct {
	Count int
}

// Copy makes a copy of a MortalityEntry
func (entry *MortalityEntry) Copy() datasets.Copyable {
	return &MortalityEntry{Count: entry.Count}
}
