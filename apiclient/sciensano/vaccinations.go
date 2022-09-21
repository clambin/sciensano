package sciensano

import (
	"fmt"
	"github.com/clambin/sciensano/apiclient"
	"time"
)

// APIVaccinationsResponse is a single entry in APIVaccinationResponse
type APIVaccinationsResponse struct {
	TimeStamp    TimeStamp `json:"DATE"`
	Manufacturer string    `json:"BRAND"`
	Region       string    `json:"REGION"`
	AgeGroup     string    `json:"AGEGROUP"`
	Gender       string    `json:"SEX"`
	Dose         DoseType  `json:"DOSE"`
	Count        int       `json:"COUNT"`
}

// APIVaccinationsResponses is a slice of APIVaccinationResponse structures
//
//easyjson:json
type APIVaccinationsResponses []*APIVaccinationsResponse

// DoseType is the type of vaccination: first, full, singledose, booster, etc.
type DoseType int

const (
	TypeVaccinationPartial DoseType = iota
	TypeVaccinationFull
	TypeVaccinationSingle
	TypeVaccinationBooster
)

func (d *DoseType) UnmarshalJSON(body []byte) (err error) {
	switch string(body) {
	case `"A"`:
		*d = TypeVaccinationPartial
	case `"B"`:
		*d = TypeVaccinationFull
	case `"C"`:
		*d = TypeVaccinationSingle
	case `"E"`, `"E2"`, `"E3"`:
		*d = TypeVaccinationBooster
	default:
		err = fmt.Errorf("invalid Dose: %s", string(body))
	}
	return
}

func (d DoseType) MarshalJSON() (body []byte, err error) {
	switch d {
	case TypeVaccinationPartial:
		body = []byte(`"A"`)
	case TypeVaccinationFull:
		body = []byte(`"B"`)
	case TypeVaccinationSingle:
		body = []byte(`"C"`)
	case TypeVaccinationBooster:
		body = []byte(`"E"`)
	default:
		err = fmt.Errorf("invalid Dose: %d", d)
	}
	return
}

var _ apiclient.APIResponse = &APIVaccinationsResponse{}

// GetTimestamp returns the entry's timestamp
func (v *APIVaccinationsResponse) GetTimestamp() time.Time {
	return v.TimeStamp.Time
}

// GetGroupFieldValue returns the value of the specified entry's field
func (v *APIVaccinationsResponse) GetGroupFieldValue(groupField int) (value string) {
	switch groupField {
	case apiclient.GroupByRegion:
		value = v.Region
	case apiclient.GroupByAgeGroup:
		value = v.AgeGroup
	case apiclient.GroupByManufacturer:
		value = v.Manufacturer
	}
	return
}

// GetTotalValue returns the entry's total number of vaccinations
func (v APIVaccinationsResponse) GetTotalValue() float64 {
	return float64(v.Count)
}

// GetAttributeNames returns the names of the types of vaccinations
func (v APIVaccinationsResponse) GetAttributeNames() []string {
	return []string{"partial", "full", "singledose", "booster"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (v APIVaccinationsResponse) GetAttributeValues() (values []float64) {
	values = []float64{0, 0, 0, 0}
	switch v.Dose {
	case TypeVaccinationPartial:
		values[0] = float64(v.Count)
	case TypeVaccinationFull:
		values[1] = float64(v.Count)
	case TypeVaccinationSingle:
		values[2] = float64(v.Count)
	case TypeVaccinationBooster:
		values[3] = float64(v.Count)
	}
	return
}
