package cases

import (
	"context"
	"fmt"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano"
	log "github.com/sirupsen/logrus"
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
		"cases-province": {TableQueryFunc: handler.buildCasesResponse},
		"cases-region":   {TableQueryFunc: handler.buildCasesResponse},
		"cases-age":      {TableQueryFunc: handler.buildCasesResponse},
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

func (handler *Handler) buildCasesResponse(ctx context.Context, target string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var cases *sciensano.Cases
	switch target {
	case "cases":
		cases, err = handler.Sciensano.GetCases(ctx, args.Range.To)
	case "cases-province":
		cases, err = handler.Sciensano.GetCasesByProvince(ctx, args.Range.To)
	case "cases-region":
		cases, err = handler.Sciensano.GetCasesByRegion(ctx, args.Range.To)
	case "cases-age":
		cases, err = handler.Sciensano.GetCasesByAgeGroup(ctx, args.Range.To)
	}

	if err != nil {
		return
	}

	timestampColumn := make(grafanajson.TableQueryResponseTimeColumn, 0, len(cases.Timestamps))
	for _, timestamp := range cases.Timestamps {
		timestampColumn = append(timestampColumn, timestamp)
	}

	response = &grafanajson.TableQueryResponse{
		Columns: []grafanajson.TableQueryResponseColumn{{
			Text: "Timestamp",
			Data: timestampColumn,
		}},
	}

	for _, group := range cases.Groups {
		dataColumn := make(grafanajson.TableQueryResponseNumberColumn, 0, len(group.Values))
		for _, value := range group.Values {
			dataColumn = append(dataColumn, float64(value))
		}

		name := group.Name
		if name == "" {
			if target == "cases" {
				name = "cases"
			} else {
				name = "(unknown)"
			}
		}

		response.Columns = append(response.Columns, grafanajson.TableQueryResponseColumn{
			Text: name,
			Data: dataColumn,
		})
	}

	return
}
