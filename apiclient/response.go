package apiclient

import "time"

type APIResponse interface {
	GetTimestamp() time.Time
	GetGroupFieldValue(groupField int) (value string)
	GetTotalValue() float64
	GetAttributeNames() []string
	GetAttributeValues() []float64
}

const (
	// GroupByAgeGroup is used by GetGroupFieldName. This groups all data by Age Group.
	GroupByAgeGroup int = iota
	// GroupByRegion is used by GetGroupFieldName. This groups all data by Region
	GroupByRegion
	// GroupByProvince is used by GetGroupFieldName. This groups all data by Province
	GroupByProvince
	// GroupByManufacturer is used by GetGroupFieldName. This groups all vaccines by manufacturer
	GroupByManufacturer
)
