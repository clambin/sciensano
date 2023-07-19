package server

import (
	"fmt"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/cache/sciensano"
)

type summaryType string

const (
	total          summaryType = "Total"
	byProvince     summaryType = "ByProvince"
	byRegion       summaryType = "ByRegion"
	byAgeGroup     summaryType = "ByAgeGroup"
	byManufacturer summaryType = "ByManufacturer"
)

var summaryTypeInfo = map[summaryType]struct {
	label  string
	column sciensano.SummaryColumn
}{
	total:          {label: "total", column: sciensano.Total},
	byProvince:     {label: "by province", column: sciensano.ByProvince},
	byRegion:       {label: "by region", column: sciensano.ByRegion},
	byAgeGroup:     {label: "by age group", column: sciensano.ByAgeGroup},
	byManufacturer: {label: "by manufacturer", column: sciensano.ByManufacturer},
}

func newMetric(name string, metricTypes ...summaryType) grafanaJSONServer.Metric {
	var options []grafanaJSONServer.MetricPayloadOption

	for _, mt := range metricTypes {
		options = append(options, grafanaJSONServer.MetricPayloadOption{
			Label: summaryTypeInfo[mt].label,
			Value: string(mt),
		})
	}
	return grafanaJSONServer.Metric{Value: name, Payloads: []grafanaJSONServer.MetricPayload{{
		Label:   "Summary",
		Name:    "summary",
		Type:    "select",
		Width:   40,
		Options: options,
	}}}
}

func getSummaryMode(target string, req grafanaJSONServer.QueryRequest) (sciensano.SummaryColumn, error) {
	var summaryOption struct {
		Summary summaryType
	}
	if err := req.GetPayload(target, &summaryOption); err != nil {
		return -1, err
	}
	mode, ok := summaryTypeInfo[summaryOption.Summary]
	if !ok {
		return -1, fmt.Errorf("invalid summary type: %s", summaryOption.Summary)
	}
	return mode.column, nil
}
