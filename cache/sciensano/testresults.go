package sciensano

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/set"
	"github.com/clambin/sciensano/pkg/tabulator"
)

type TestResult struct {
	TimeStamp TimeStamp `json:"DATE"`
	Province  string    `json:"PROVINCE"`
	Region    string    `json:"REGION"`
	Total     int       `json:"TESTS_ALL"`
	Positive  int       `json:"TESTS_ALL_POS"`
}

type TestResults []*TestResult

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
		if columnName == "" {
			columnName = "(unknown)"
		}
		if columnNames.IsNew(columnName) {
			t.RegisterColumn(columnName)
		}

		t.Add(testResult.TimeStamp.Time, columnName, float64(testResult.Total))
	}

	return t, nil
}

func (r TestResults) Categorize() *tabulator.Tabulator {
	t := tabulator.New("positive", "total", "rate")

	for _, testResult := range r {
		t.Add(testResult.TimeStamp.Time, "positive", float64(testResult.Positive))
		t.Add(testResult.TimeStamp.Time, "total", float64(testResult.Total))
	}

	positive, _ := t.GetValues("positive")
	total, _ := t.GetValues("total")

	for index, timestamp := range t.GetTimestamps() {
		t.Add(timestamp, "rate", positive[index]/total[index])
	}

	return t
}
