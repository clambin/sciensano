package sciensano

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/set"
	"github.com/clambin/sciensano/pkg/tabulator"
)

type Hospitalisation struct {
	TimeStamp   TimeStamp `json:"DATE"`
	Province    string    `json:"PROVINCE"`
	Region      string    `json:"REGION"`
	TotalIn     int       `json:"TOTAL_IN"`
	TotalInICU  int       `json:"TOTAL_IN_ICU"`
	TotalInResp int       `json:"TOTAL_IN_RESP"`
	TotalInECMO int       `json:"TOTAL_IN_ECMO"`
}

type Hospitalisations []*Hospitalisation

func (h Hospitalisations) Summarize(summaryColumn SummaryColumn) (*tabulator.Tabulator, error) {
	t := tabulator.New()

	var columnNames set.Set
	for _, hospitalisation := range h {
		var columnName string
		switch summaryColumn {
		case Total:
			columnName = "Total"
		case ByRegion:
			columnName = hospitalisation.Region
		case ByProvince:
			columnName = hospitalisation.Province
		default:
			return nil, fmt.Errorf("hospitalisations: invalid summary column: %s", summaryColumn.String())
		}
		if columnName == "" {
			columnName = "(unknown)"
		}
		if columnNames.IsNew(columnName) {
			t.RegisterColumn(columnName)
		}

		t.Add(hospitalisation.TimeStamp.Time, columnName, float64(hospitalisation.TotalIn))
	}

	return t, nil
}

func (h Hospitalisations) Categorize() *tabulator.Tabulator {
	t := tabulator.New("in", "inICU", "inResp", "inECMO")

	for _, hospitalisation := range h {
		t.Add(hospitalisation.TimeStamp.Time, "in", float64(hospitalisation.TotalIn))
		t.Add(hospitalisation.TimeStamp.Time, "inICU", float64(hospitalisation.TotalInICU))
		t.Add(hospitalisation.TimeStamp.Time, "inResp", float64(hospitalisation.TotalInResp))
		t.Add(hospitalisation.TimeStamp.Time, "inECMO", float64(hospitalisation.TotalInECMO))
	}

	return t
}
