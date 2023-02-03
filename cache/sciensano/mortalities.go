package sciensano

import (
	"fmt"
	"github.com/clambin/go-common/set"
	"github.com/clambin/go-common/tabulator"
)

type Mortality struct {
	TimeStamp TimeStamp `json:"DATE"`
	Region    string    `json:"REGION"`
	AgeGroup  string    `json:"AGEGROUP"`
	Deaths    int       `json:"DEATHS"`
}

type Mortalities []*Mortality

func (m Mortalities) Summarize(summaryColumn SummaryColumn) (*tabulator.Tabulator, error) {
	t := tabulator.New()

	columnNames := set.Create[string]()
	for _, mortality := range m {
		var columnName string
		switch summaryColumn {
		case Total:
			columnName = "Total"
		case ByRegion:
			columnName = mortality.Region
		case ByAgeGroup:
			columnName = mortality.AgeGroup
		default:
			return nil, fmt.Errorf("mortalities: invalid summary column: %s", summaryColumn.String())
		}
		if columnName == "" {
			columnName = "(unknown)"
		}
		if !columnNames.Contains(columnName) {
			t.RegisterColumn(columnName)
			columnNames.Add(columnName)
		}

		t.Add(mortality.TimeStamp.Time, columnName, float64(mortality.Deaths))
	}

	return t, nil
}
