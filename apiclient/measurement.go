package apiclient

import "time"

type Measurement interface {
	GetTimestamp() time.Time
	GetGroupFieldValue(groupField int) string
}

type Measurements []Measurement

const (
	GroupByNone int = iota
	GroupByAgeGroup
	GroupByRegion
	GroupByProvince
)

var _ Measurement = &APICasesResponseEntry{}
var _ Measurement = &APIHospitalisationsResponseEntry{}
var _ Measurement = &APIMortalityResponseEntry{}
var _ Measurement = &APITestResultsResponseEntry{}
var _ Measurement = &APIVaccinationsResponseEntry{}
