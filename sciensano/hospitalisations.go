package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/datasets"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

// HospitalisationsGetter contains all methods providing COVID-19-related hospitalisation figures
type HospitalisationsGetter interface {
	GetHospitalisations(ctx context.Context) (results *datasets.Dataset, err error)
	GetHospitalisationsByRegion(ctx context.Context) (results *datasets.Dataset, err error)
	GetHospitalisationsByProvince(ctx context.Context) (results *datasets.Dataset, err error)
}

const (
	groupHospitalisationsByNone = iota
	groupHospitalisationsByRegion
	groupHospitalisationsByProvince
)

// GetHospitalisations returns all hospitalisations
func (client *Client) GetHospitalisations(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getHospitalisations(ctx, "GetHospitalisations", "Hospitalisations", groupHospitalisationsByNone)
}

// GetHospitalisationsByRegion returns all hospitalisations, grouped by region
func (client *Client) GetHospitalisationsByRegion(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getHospitalisations(ctx, "GetHospitalisationsByRegion", "HospitalisationsByRegion", groupHospitalisationsByRegion)
}

// GetHospitalisationsByProvince returns all hospitalisations, grouped by province
func (client *Client) GetHospitalisationsByProvince(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getHospitalisations(ctx, "GetHospitalisationsByProvince", "HospitalisationsByProvince", groupHospitalisationsByProvince)
}

func (client *Client) getHospitalisations(ctx context.Context, name, cacheEntryName string, mode int) (results *datasets.Dataset, err error) {
	before := time.Now()
	defer func() { log.WithField("time", time.Now().Sub(before)).Debug(name + " done") }()

	log.Debug("running " + name)
	entry := client.cache.Load(cacheEntryName)
	entry.Once.Do(func() {
		var apiResult apiclient.APIHospitalisationsResponse
		if apiResult, err = client.Getter.GetHospitalisations(ctx); err == nil {
			entry.Data = groupHospitalisations(apiResult, mode)
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

func groupHospitalisations(hospitalisations apiclient.APIHospitalisationsResponse, groupField int) (results *datasets.Dataset) {
	mappedHospitalisations, mappedGroups := mapHospitalisations(hospitalisations, groupField)
	timestamps := getUniqueSortedTimestampsFromHospitalisations(mappedHospitalisations)
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
			value, ok := mappedHospitalisations[timestamp][group]
			if ok == false {
				value = &HospitalisationsEntry{}
			}
			results.Groups[index].Values = append(results.Groups[index].Values, value)
		}
	}
	return

}

func mapHospitalisations(hospitalisations apiclient.APIHospitalisationsResponse, groupField int) (mappedHospitalisations map[time.Time]map[string]*HospitalisationsEntry, mappedGroups map[string]struct{}) {
	mappedHospitalisations = make(map[time.Time]map[string]*HospitalisationsEntry)
	mappedGroups = make(map[string]struct{})

	for _, entry := range hospitalisations {
		mappedHospitalisation, ok := mappedHospitalisations[entry.TimeStamp.Time]
		if ok == false {
			mappedHospitalisation = make(map[string]*HospitalisationsEntry)
		}

		var groupName string
		switch groupField {
		case groupHospitalisationsByNone:
			groupName = ""
		case groupHospitalisationsByRegion:
			groupName = entry.Region
		case groupHospitalisationsByProvince:
			groupName = entry.Province
		}

		value, ok := mappedHospitalisation[groupName]
		if ok == false {
			value = &HospitalisationsEntry{}
		}

		value.Add(entry)
		mappedHospitalisation[groupName] = value
		mappedHospitalisations[entry.TimeStamp.Time] = mappedHospitalisation

		mappedGroups[groupName] = struct{}{}
	}

	return
}

func getUniqueSortedTimestampsFromHospitalisations(input map[time.Time]map[string]*HospitalisationsEntry) (output []time.Time) {
	for timestamp := range input {
		output = append(output, timestamp)
	}
	sort.Slice(output, func(i, j int) bool { return output[i].Before(output[j]) })
	return
}
