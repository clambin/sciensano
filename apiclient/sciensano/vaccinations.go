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
	TypeVaccinationBooster2
	TypeVaccinationBooster3
)

func (d *DoseType) UnmarshalJSON(body []byte) (err error) {
	switch string(body) {
	case `"A"`:
		*d = TypeVaccinationPartial
	case `"B"`:
		*d = TypeVaccinationFull
	case `"C"`:
		*d = TypeVaccinationSingle
	case `"E"`:
		*d = TypeVaccinationBooster
	case `"E2"`:
		*d = TypeVaccinationBooster2
	case `"E3"`:
		*d = TypeVaccinationBooster3
	default:
		err = fmt.Errorf("invalid Dose: %s", string(body))
	}
	return
}

func (d *DoseType) MarshalJSON() (body []byte, err error) {
	switch *d {
	case TypeVaccinationPartial:
		body = []byte(`"A"`)
	case TypeVaccinationFull:
		body = []byte(`"B"`)
	case TypeVaccinationSingle:
		body = []byte(`"C"`)
	case TypeVaccinationBooster:
		body = []byte(`"E"`)
	case TypeVaccinationBooster2:
		body = []byte(`"E2"`)
	case TypeVaccinationBooster3:
		body = []byte(`"E3"`)
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
func (v *APIVaccinationsResponse) GetTotalValue() float64 {
	return float64(v.Count)
}

// GetAttributeNames returns the names of the types of vaccinations
func (v *APIVaccinationsResponse) GetAttributeNames() []string {
	return []string{"partial", "full", "singledose", "booster", "booster2", "booster3"}
}

// GetAttributeValues gets the value for each supported type of vaccination
func (v *APIVaccinationsResponse) GetAttributeValues() (values []float64) {
	values = make([]float64, len(v.GetAttributeNames()))
	switch v.Dose {
	case TypeVaccinationPartial:
		values[0] = float64(v.Count)
	case TypeVaccinationFull:
		values[1] = float64(v.Count)
	case TypeVaccinationSingle:
		values[2] = float64(v.Count)
	case TypeVaccinationBooster:
		values[3] = float64(v.Count)
	case TypeVaccinationBooster2:
		values[4] = float64(v.Count)
	case TypeVaccinationBooster3:
		values[5] = float64(v.Count)
	}
	return
}
