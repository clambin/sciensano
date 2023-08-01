package sciensano

import (
	"fmt"
	"github.com/clambin/go-common/set"
	"github.com/clambin/go-common/tabulator"
)

//easyjson:json
type Case struct {
	TimeStamp TimeStamp `json:"DATE"`
	Province  string    `json:"PROVINCE"`
	Region    string    `json:"REGION"`
	AgeGroup  string    `json:"AGEGROUP"`
	Cases     int       `json:"CASES"`
}

//easyjson:json
type Cases []*Case

func CasesValidSummaryModes() set.Set[SummaryColumn] {
	return set.Create(Total, ByRegion, ByProvince, ByAgeGroup)
}

func (cs Cases) Summarize(summaryColumn SummaryColumn) (*tabulator.Tabulator, error) {
	t := tabulator.New()

	columnNames := set.Create[string]()
	for _, c := range cs {
		var columnName string
		switch summaryColumn {
		case Total:
			columnName = "Total"
		case ByRegion:
			columnName = c.Region
		case ByProvince:
			columnName = c.Province
		case ByAgeGroup:
			columnName = c.AgeGroup
		default:
			return nil, fmt.Errorf("cases: invalid summary column: %s", summaryColumn.String())
		}
		if columnName == "" {
			columnName = "(unknown)"
		}
		if !columnNames.Contains(columnName) {
			t.RegisterColumn(columnName)
			columnNames.Add(columnName)
		}

		t.Add(c.TimeStamp.Time, columnName, float64(c.Cases))
	}

	return t, nil
}
