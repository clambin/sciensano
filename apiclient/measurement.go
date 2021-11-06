package apiclient

import "time"

// Measurement represent a data measurement returned by one of the Sciensano APIs.
type Measurement interface {
	GetTimestamp() time.Time
	GetGroupFieldValue(groupField int) string
}

const (
	// GroupByNone is used by GetGroupFieldValue. This groups all data regardless of age, region, province, etc.
	GroupByNone int = iota
	// GroupByAgeGroup is used by GetGroupFieldValue. This groups all data by Age Group.
	GroupByAgeGroup
	// GroupByRegion is used by GetGroupFieldValue. This groups all data by Region
	GroupByRegion
	// GroupByProvince is used by GetGroupFieldValue. This groups all data by Province
	GroupByProvince
)

var _ Measurement = &APICasesResponseEntry{}
var _ Measurement = &APIHospitalisationsResponseEntry{}
var _ Measurement = &APIMortalityResponseEntry{}
var _ Measurement = &APITestResultsResponseEntry{}
var _ Measurement = &APIVaccinationsResponseEntry{}
