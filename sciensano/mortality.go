package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/datasets"
	log "github.com/sirupsen/logrus"
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

// GetMortality returns all mortality figures
func (client *Client) GetMortality(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getMortality(ctx, "GetMortality", "Mortality", apiclient.GroupByNone)
}

// GetMortalityByRegion returns all mortality figures, grouped by region
func (client *Client) GetMortalityByRegion(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getMortality(ctx, "GetMortalityByRegion", "MortalityByRegion", apiclient.GroupByRegion)
}

// GetMortalityByAgeGroup returns all Mortality, grouped by age group
func (client *Client) GetMortalityByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getMortality(ctx, "GetMortalityByAge", "MortalityByAge", apiclient.GroupByAgeGroup)
}

func (client *Client) getMortality(ctx context.Context, name, cacheEntryName string, mode int) (results *datasets.Dataset, err error) {
	before := time.Now()
	defer func() { log.WithField("time", time.Now().Sub(before)).Debug(name + " done") }()

	log.Debug("running " + name)
	entry := client.Cache.Load(cacheEntryName)
	entry.Once.Do(func() {
		var apiResult []apiclient.Measurement
		if apiResult, err = client.Getter.GetMortality(ctx); err == nil {
			entry.Data = groupMeasurements(apiResult, mode, NewMortalityEntry)
			client.Cache.Save(cacheEntryName, entry)
		} else {
			client.Cache.Clear(cacheEntryName)
		}
	})
	if err == nil && entry.Data != nil {
		results = entry.Data.Copy()
	}
	return
}

// MortalityEntry contains the mortality for a single timestamp
type MortalityEntry struct {
	Count int
}

// NewMortalityEntry returns a new MortalityEntry, as a GroupedEntry. Used by groupMeasurements
func NewMortalityEntry() GroupedEntry {
	return &MortalityEntry{}
}

// Copy makes a copy of a MortalityEntry
func (entry *MortalityEntry) Copy() datasets.Copyable {
	return &MortalityEntry{Count: entry.Count}
}

// Add adds the passed MortalityEntry values to its own values
func (entry *MortalityEntry) Add(input apiclient.Measurement) {
	entry.Count += input.(*apiclient.APIMortalityResponseEntry).Deaths
}
