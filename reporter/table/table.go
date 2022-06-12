package table

import (
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/reporter/table/tabulator"
	"github.com/clambin/simplejson/v3/data"
	grafanaData "github.com/grafana/grafana-plugin-sdk-go/data"
	"time"
)

// NewFromAPIResponse creates a new Dataset from a slice of API responses. API responses need to be in sequential order.
func NewFromAPIResponse(response []apiclient.APIResponse) *data.Table {
	rows := len(response)
	timestamps := make([]time.Time, 0, rows)

	var values [][]float64
	var attribs []string
	var lastTimestamp time.Time
	row := -1

	for _, entry := range response {
		// first entry. create a column for each attribute name
		if lastTimestamp.IsZero() {
			attribs = entry.GetAttributeNames()
			for range attribs {
				values = append(values, make([]float64, 0, rows))
			}
		}
		// new timestamp. Create a new row.
		if !entry.GetTimestamp().Equal(lastTimestamp) {
			row++
			lastTimestamp = entry.GetTimestamp()
			timestamps = append(timestamps, lastTimestamp)
			for i := range attribs {
				values[i] = append(values[i], 0.0)
			}
		}
		// add new values to the current row
		for i, value := range entry.GetAttributeValues() {
			values[i][row] += value
		}
	}

	fields := make(grafanaData.Fields, 1+len(values))
	fields[0] = grafanaData.NewField("time", nil, timestamps)
	for idx, attrib := range attribs {
		fields[idx+1] = grafanaData.NewField(attrib, nil, values[idx])
	}
	return &data.Table{Frame: grafanaData.NewFrame("frame", fields...)}
}

// NewGroupedFromAPIResponse created a (grouped) dataframe from a set of APIResponses. Responses need to be in sequential order.
func NewGroupedFromAPIResponse(response []apiclient.APIResponse, groupField int) *data.Table {
	table := tabulator.New(getUniqueColumns(response, groupField)...)

	for _, entry := range response {
		_ = table.Add(entry.GetTimestamp(), entry.GetGroupFieldValue(groupField), entry.GetTotalValue())
	}

	return &data.Table{Frame: tableToDataFrame(table)}
}

func getUniqueColumns(response []apiclient.APIResponse, groupField int) (columns []string) {
	cols := make(map[string]struct{})
	for _, entry := range response {
		cols[entry.GetGroupFieldValue(groupField)] = struct{}{}
	}
	for col := range cols {
		columns = append(columns, col)
	}
	return
}

func tableToDataFrame(table *tabulator.Tabulator) *grafanaData.Frame {
	fields := make(grafanaData.Fields, 0, table.Size())
	fields = append(fields, grafanaData.NewField("time", nil, table.GetTimestamps()))
	for _, col := range table.GetColumns() {
		values, _ := table.GetValues(col)
		fields = append(fields, grafanaData.NewField(col, nil, values))
	}
	return grafanaData.NewFrame("frame", fields...)
}
