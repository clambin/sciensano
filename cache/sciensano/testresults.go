package sciensano

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/set"
	"github.com/clambin/sciensano/reporter/table/tabulator"
)

type TestResult struct {
	TimeStamp TimeStamp `json:"DATE"`
	Province  string    `json:"PROVINCE"`
	Region    string    `json:"REGION"`
	Total     int       `json:"TESTS_ALL"`
	Positive  int       `json:"TESTS_ALL_POS"`
}

type TestResults []TestResult

func (r TestResults) Summarize(summaryColumn SummaryColumn) (*tabulator.Tabulator, error) {
	t := tabulator.New()

	var columnNames set.Set
	for _, testResult := range r {
		var columnName string
		switch summaryColumn {
		case Total:
			columnName = "Total"
		case ByRegion:
			columnName = testResult.Region
		case ByProvince:
			columnName = testResult.Province
		default:
			return nil, fmt.Errorf("testResults: invalid summary column: %s", summaryColumn.String())
		}
		if columnNames.IsNew(columnName) {
			t.RegisterColumn(columnName)
		}

		t.Add(testResult.TimeStamp.Time, columnName, float64(testResult.Total))
	}

	return t, nil
}
