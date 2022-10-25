package apiclient

import "time"

// APIResponse is the interface that each apiclient response needs to adhere to.
// These allow datasets to process API responses in a generic way
type APIResponse interface {
	GetTimestamp() time.Time
	GetGroupFieldValue(groupField GroupField) (value string)
	GetTotalValue() float64
	GetAttributeNames() []string
	GetAttributeValues() []float64
}

type GroupField int

const (
	// GroupByAgeGroup is used by GetGroupFieldName. This groups all data by Age Group.
	GroupByAgeGroup GroupField = iota
	// GroupByRegion is used by GetGroupFieldName. This groups all data by Region
	GroupByRegion
	// GroupByProvince is used by GetGroupFieldName. This groups all data by Province
	GroupByProvince
	// GroupByManufacturer is used by GetGroupFieldName. This groups all vaccines by manufacturer
	GroupByManufacturer
)
