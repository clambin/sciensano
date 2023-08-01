package sciensano

import (
	"fmt"
	"github.com/clambin/go-common/set"
	"github.com/clambin/go-common/tabulator"
)

// DoseType is the type of vaccination: first, full, singledose, booster, etc.
type DoseType int

const (
	Partial DoseType = iota
	Full
	SingleDose
	Booster
	Booster2
	Booster3
)

var doseTypeStrings = map[DoseType]string{
	Partial:    "Partial",
	Full:       "Full",
	SingleDose: "SingleDose",
	Booster:    "Booster",
	Booster2:   "Booster 2",
	Booster3:   "Booster 3",
}

func (d *DoseType) String() string {
	value, ok := doseTypeStrings[*d]
	if !ok {
		value = "(unknown)"
	}
	return value
}

func (d *DoseType) UnmarshalJSON(body []byte) (err error) {
	switch string(body) {
	case `"A"`:
		*d = Partial
	case `"B"`:
		*d = Full
	case `"C"`:
		*d = SingleDose
	case `"E"`:
		*d = Booster
	case `"E2"`:
		*d = Booster2
	case `"E3"`:
		*d = Booster3
	default:
		err = fmt.Errorf("invalid Dose: %s", string(body))
	}
	return
}

func (d DoseType) MarshalJSON() (body []byte, err error) {
	switch d {
	case Partial:
		body = []byte(`"A"`)
	case Full:
		body = []byte(`"B"`)
	case SingleDose:
		body = []byte(`"C"`)
	case Booster:
		body = []byte(`"E"`)
	case Booster2:
		body = []byte(`"E2"`)
	case Booster3:
		body = []byte(`"E3"`)
	default:
		err = fmt.Errorf("invalid Dose: %d", d)
	}
	return
}

//easyjson:json
type Vaccination struct {
	TimeStamp    TimeStamp `json:"DATE"`
	Manufacturer string    `json:"BRAND"`
	Region       string    `json:"REGION"`
	AgeGroup     string    `json:"AGEGROUP"`
	Gender       string    `json:"SEX"`
	Dose         DoseType  `json:"DOSE"`
	Count        int       `json:"COUNT"`
}

//easyjson:json
type Vaccinations []*Vaccination

func VaccinationsValidSummaryModes() set.Set[SummaryColumn] {
	return set.Create(Total, ByRegion, ByAgeGroup, ByManufacturer, ByVaccinationType)
}

func (v Vaccinations) Summarize(summaryColumn SummaryColumn) (*tabulator.Tabulator, error) {
	t := tabulator.New()

	columnNames := set.Create[string]()
	for _, vaccination := range v {
		var columnName string
		switch summaryColumn {
		case Total:
			columnName = "Total"
		case ByRegion:
			columnName = vaccination.Region
		case ByAgeGroup:
			columnName = vaccination.AgeGroup
		case ByManufacturer:
			columnName = vaccination.Manufacturer
		case ByVaccinationType:
			switch vaccination.Dose {
			case Partial:
				columnName = "partial"
			case Full, SingleDose:
				columnName = "full"
			case Booster:
				columnName = "booster"
			case Booster2:
				columnName = "booster2"
			case Booster3:
				columnName = "booster3"
			}
		default:
			return nil, fmt.Errorf("vaccinations: invalid summary column: %s", summaryColumn.String())
		}
		if columnName == "" {
			columnName = "(unknown)"
		}
		if !columnNames.Contains(columnName) {
			t.RegisterColumn(columnName)
			columnNames.Add(columnName)
		}

		t.Add(vaccination.TimeStamp.Time, columnName, float64(vaccination.Count))
	}

	return t, nil
}
