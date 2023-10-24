package sciensano

import (
	"fmt"
	"github.com/clambin/go-common/set"
	"github.com/clambin/go-common/tabulator"
)

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

func (v Vaccination) GetSummaryColumnName(column SummaryColumn) (string, error) {
	var columnName string
	switch column {
	case Total:
		columnName = "Total"
	case ByRegion:
		columnName = v.Region
	case ByAgeGroup:
		columnName = v.AgeGroup
	case ByManufacturer:
		columnName = v.Manufacturer
	case ByVaccinationType:
		columnName = v.Dose.String()
	default:
		return "nil", fmt.Errorf("invalid summary column: %s", column.String())
	}
	return columnName, nil
}

//easyjson:json
type Vaccinations []Vaccination

func VaccinationsValidSummaryModes() set.Set[SummaryColumn] {
	return set.Create(Total, ByRegion, ByAgeGroup, ByManufacturer, ByVaccinationType)
}

func (v Vaccinations) Summarize(summaryColumn SummaryColumn) (*tabulator.Tabulator, error) {
	t := tabulator.New()

	columnNames := set.Create[string]()
	for _, vaccination := range v {
		columnName, err := vaccination.GetSummaryColumnName(summaryColumn)
		if err != nil {
			return nil, fmt.Errorf("summary: %w", err)
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
