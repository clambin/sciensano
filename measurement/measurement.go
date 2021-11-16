package measurement

import (
	"time"
)

// Measurement represent a data measurement returned by one of the Reporter APIs.
type Measurement interface {
	GetTimestamp() time.Time
	GetGroupFieldValue(groupField int) string
	GetTotalValue() float64
	GetAttributeNames() []string
	GetAttributeValues() []float64
}

const (
	// GroupByNone is used by GetGroupFieldName. This groups all data regardless of age, region, province, etc.
	GroupByNone int = iota
	// GroupByAgeGroup is used by GetGroupFieldName. This groups all data by Age Group.
	GroupByAgeGroup
	// GroupByRegion is used by GetGroupFieldName. This groups all data by Region
	GroupByRegion
	// GroupByProvince is used by GetGroupFieldName. This groups all data by Province
	GroupByProvince
	// GroupByManufacturer is used by GetGroupFieldNane. This groups all vaccines by manufacturer
	GroupByManufacturer
)
