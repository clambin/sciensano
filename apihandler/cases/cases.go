package cases

import (
	"context"
	"fmt"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

// Handler implements a grafana-json handler for COVID-19 cases
type Handler struct {
	apiclient.Getter
	Sciensano   sciensano.APIClient
	targetTable grafanajson.TargetTable
}

// New creates a new Handler
func New(getter apiclient.Getter, client sciensano.APIClient) (handler *Handler) {
	handler = &Handler{
		Getter:    getter,
		Sciensano: client,
	}

	handler.targetTable = grafanajson.TargetTable{
		"cases":          {TableQueryFunc: handler.buildCasesResponse},
		"cases-province": {TableQueryFunc: handler.buildGroupedCases},
		"cases-region":   {TableQueryFunc: handler.buildGroupedCases},
	}

	return
}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (handler *Handler) Endpoints() grafanajson.Endpoints {
	return grafanajson.Endpoints{
		Search:     handler.Search,
		TableQuery: handler.TableQuery,
	}
}

// Search implements the grafana-json Search function. It returns all supported targets
func (handler *Handler) Search() (targets []string) {
	return handler.targetTable.Targets()
}

// TableQuery implements the grafana-json TableQuery function. It processes incoming TableQuery requests
func (handler *Handler) TableQuery(ctx context.Context, target string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	start := time.Now()
	response, err = handler.targetTable.RunTableQuery(ctx, target, args)
	if err != nil {
		return nil, fmt.Errorf("unable to build table query response for target '%s': %s", target, err.Error())
	}
	log.WithFields(log.Fields{"duration": time.Now().Sub(start), "target": target}).Debug("TableQuery called")
	return
}

func (handler *Handler) buildCasesResponse(ctx context.Context, _ string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var cases []sciensano.CaseCount
	cases, err = handler.Sciensano.GetCases(ctx, args.Range.To)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve cases: %s", err.Error())
	}

	rows := len(cases)
	timestamps := make(grafanajson.TableQueryResponseTimeColumn, rows)
	casesData := make(grafanajson.TableQueryResponseNumberColumn, rows)

	for index, entry := range cases {
		timestamps[index] = entry.Timestamp
		casesData[index] = float64(entry.Count)
	}

	response = new(grafanajson.TableQueryResponse)
	response.Columns = []grafanajson.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestamps},
		{Text: "cases", Data: casesData},
	}

	return
}

func (handler *Handler) buildGroupedCases(ctx context.Context, target string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var cases map[string][]sciensano.CaseCount
	switch target {
	case "cases-province":
		cases, err = handler.Sciensano.GetCasesByProvince(ctx, args.Range.To)
	case "cases-region":
		cases, err = handler.Sciensano.GetCasesByRegion(ctx, args.Range.To)
	}

	if err != nil {
		return
	}

	// TODO: more efficient way to get timestamps & keys while building groupedCases?
	groupedCases := fillCases(cases)
	timestamps := getTimestamps(groupedCases)
	keys := getKeys(cases)

	timestampColumn := make(grafanajson.TableQueryResponseTimeColumn, 0, len(timestamps))
	for _, timestamp := range timestamps {
		timestampColumn = append(timestampColumn, timestamp)
	}

	response = &grafanajson.TableQueryResponse{
		Columns: []grafanajson.TableQueryResponseColumn{{
			Text: "Timestamp",
			Data: timestampColumn,
		}},
	}

	for _, key := range keys {
		dataColumn := make(grafanajson.TableQueryResponseNumberColumn, 0, len(timestamps))
		for _, timestamp := range timestamps {
			dataColumn = append(dataColumn, float64(groupedCases[timestamp][key]))
		}

		if key == "" {
			key = "(unknown)"
		}

		response.Columns = append(response.Columns, grafanajson.TableQueryResponseColumn{
			Text: key,
			Data: dataColumn,
		})
	}

	return
}

func fillCases(cases map[string][]sciensano.CaseCount) (results map[time.Time]map[string]int) {
	results = make(map[time.Time]map[string]int)

	for key, data := range cases {
		for _, entry := range data {
			_, ok := results[entry.Timestamp]
			if ok == false {
				results[entry.Timestamp] = make(map[string]int)
				for key2 := range cases {
					results[entry.Timestamp][key2] = 0
				}
			}
			count := results[entry.Timestamp][key]
			count += entry.Count
			results[entry.Timestamp][key] = count
		}
	}
	return
}

func getKeys(cases map[string][]sciensano.CaseCount) (results []string) {
	for key := range cases {
		results = append(results, key)
	}
	sort.Strings(results)
	return
}

func getTimestamps(groupedCases map[time.Time]map[string]int) (timestamps []time.Time) {
	for timestamp := range groupedCases {
		timestamps = append(timestamps, timestamp)
	}
	sort.Slice(timestamps, func(i, j int) bool { return timestamps[i].Before(timestamps[j]) })
	return
}
