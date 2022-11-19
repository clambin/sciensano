package sciensano

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/set"
	"github.com/clambin/sciensano/reporter/table/tabulator"
)

type Case struct {
	TimeStamp TimeStamp `json:"DATE"`
	Province  string    `json:"PROVINCE"`
	Region    string    `json:"REGION"`
	AgeGroup  string    `json:"AGEGROUP"`
	Cases     int       `json:"CASES"`
}

type Cases []Case

func (cs Cases) Summarize(summaryColumn SummaryColumn) (*tabulator.Tabulator, error) {
	t := tabulator.New()

	var columnNames set.Set
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
		if columnNames.IsNew(columnName) {
			t.RegisterColumn(columnName)
		}

		t.Add(c.TimeStamp.Time, columnName, float64(c.Cases))
	}

	return t, nil
}
