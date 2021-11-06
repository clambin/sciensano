package sciensano

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/datasets"
	log "github.com/sirupsen/logrus"
	"time"
)

// HospitalisationsGetter contains all methods providing COVID-19-related hospitalisation figures
type HospitalisationsGetter interface {
	GetHospitalisations(ctx context.Context) (results *datasets.Dataset, err error)
	GetHospitalisationsByRegion(ctx context.Context) (results *datasets.Dataset, err error)
	GetHospitalisationsByProvince(ctx context.Context) (results *datasets.Dataset, err error)
}

// GetHospitalisations returns all hospitalisations
func (client *Client) GetHospitalisations(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getHospitalisations(ctx, "GetHospitalisations", "Hospitalisations", apiclient.GroupByNone)
}

// GetHospitalisationsByRegion returns all hospitalisations, grouped by region
func (client *Client) GetHospitalisationsByRegion(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getHospitalisations(ctx, "GetHospitalisationsByRegion", "HospitalisationsByRegion", apiclient.GroupByRegion)
}

// GetHospitalisationsByProvince returns all hospitalisations, grouped by province
func (client *Client) GetHospitalisationsByProvince(ctx context.Context) (results *datasets.Dataset, err error) {
	return client.getHospitalisations(ctx, "GetHospitalisationsByProvince", "HospitalisationsByProvince", apiclient.GroupByProvince)
}

func (client *Client) getHospitalisations(ctx context.Context, name, cacheEntryName string, mode int) (results *datasets.Dataset, err error) {
	before := time.Now()
	defer func() { log.WithField("time", time.Now().Sub(before)).Debug(name + " done") }()

	log.Debug("running " + name)
	entry := client.Cache.Load(cacheEntryName)
	entry.Once.Do(func() {
		var apiResult []apiclient.Measurement
		if apiResult, err = client.Getter.GetHospitalisations(ctx); err == nil {
			entry.Data = groupMeasurements(apiResult, mode, NewHospitalisationsEntry)
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
