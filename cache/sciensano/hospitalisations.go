package sciensano

import (
	"fmt"
	"github.com/clambin/sciensano/pkg/set"
	"github.com/clambin/sciensano/reporter/table/tabulator"
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

type Hospitalisations []Hospitalisation

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
		if columnNames.IsNew(columnName) {
			t.RegisterColumn(columnName)
		}

		t.Add(hospitalisation.TimeStamp.Time, columnName, float64(hospitalisation.TotalIn))
	}

	return t, nil
}
