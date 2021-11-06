package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/datasets"
	log "github.com/sirupsen/logrus"
	"time"
)

// CasesGetter contains all methods providing COVID-19 cases
type CasesGetter interface {
	GetCases(ctx context.Context) (results *datasets.Dataset, err error)
	GetCasesByRegion(ctx context.Context) (results *datasets.Dataset, err error)
	GetCasesByProvince(ctx context.Context) (results *datasets.Dataset, err error)
	GetCasesByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error)
}

// GetCases returns all cases
func (client *Client) GetCases(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getCases(ctx, "GetCases", "Cases", apiclient.GroupByNone)
}

// GetCasesByRegion returns all cases, grouped by region
func (client *Client) GetCasesByRegion(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getCases(ctx, "GetCasesByRegion", "CasesByRegion", apiclient.GroupByRegion)
}

// GetCasesByProvince returns all cases, grouped by province
func (client *Client) GetCasesByProvince(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getCases(ctx, "GetCasesByProvince", "CasesByProvince", apiclient.GroupByProvince)
}

// GetCasesByAgeGroup returns all cases, grouped by province
func (client *Client) GetCasesByAgeGroup(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getCases(ctx, "GetCasesByAgeGroup", "CasesByAgeGroup", apiclient.GroupByAgeGroup)
}

// getCases
func (client *Client) getCases(ctx context.Context, name, cacheEntryName string, mode int) (results *datasets.Dataset, err error) {
	before := time.Now()
	defer func() { log.WithField("time", time.Now().Sub(before)).Debug(name + " done") }()

	log.Debug("running " + name)
	entry := client.Cache.Load(cacheEntryName)
	entry.Once.Do(func() {
		var apiResult []apiclient.Measurement
		if apiResult, err = client.Getter.GetCases(ctx); err == nil {
			entry.Data = groupMeasurements(apiResult, mode, NewCasesEntry)
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
