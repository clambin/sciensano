package sciensano

import (
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/datasets"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

type GroupedEntry interface {
	datasets.Copyable
	Add(entry apiclient.Measurement)
}

func groupMeasurements(entries []apiclient.Measurement, groupField int, newEntry func() GroupedEntry) (results *datasets.Dataset) {
	before := time.Now()
	defer log.WithField("time", time.Now().Sub(before)).Debug("groupMortality")

	mappedMortalities, mappedGroups := mapEntries(entries, groupField, newEntry)
	timestamps := getUniqueSortedTimestamps(mappedMortalities)
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
				entry = newEntry()
			}
			results.Groups[index].Values = append(results.Groups[index].Values, entry)
		}
	}
	return
}

func mapEntries(entries []apiclient.Measurement, groupField int, newEntryFunc func() GroupedEntry) (mappedEntries map[time.Time]map[string]GroupedEntry, mappedGroups map[string]struct{}) {
	before := time.Now()
	defer log.WithField("time", time.Now().Sub(before)).Debug("mapMortality")

	mappedEntries = make(map[time.Time]map[string]GroupedEntry)
	mappedGroups = make(map[string]struct{})

	currentTimestamp := time.Time{}
	currentEntries := make(map[string]GroupedEntry)

	for _, newEntry := range entries {
		entryTimestamp := newEntry.(apiclient.Measurement).GetTimestamp()
		groupName := newEntry.(apiclient.Measurement).GetGroupFieldValue(groupField)

		if !currentTimestamp.IsZero() && !currentTimestamp.Equal(entryTimestamp) {
			mappedEntries[currentTimestamp] = currentEntries
			currentEntries = make(map[string]GroupedEntry)
		}
		currentTimestamp = entryTimestamp

		mappedGroups[groupName] = struct{}{}

		entry, ok := currentEntries[groupName]

		if ok == false {
			entry = newEntryFunc()
		}

		entry.Add(newEntry.(apiclient.Measurement))
		currentEntries[groupName] = entry
	}

	if !currentTimestamp.IsZero() {
		mappedEntries[currentTimestamp] = currentEntries
	}
	return
}

func getUniqueSortedTimestamps(input map[time.Time]map[string]GroupedEntry) (output []time.Time) {
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
