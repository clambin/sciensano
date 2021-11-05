package hospitalisations

import (
	"context"
	"fmt"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/sciensano/datasets"
	log "github.com/sirupsen/logrus"
	"time"
)

// Handler implements a grafana-json handler for COVID-19 cases
type Handler struct {
	Sciensano   sciensano.APIClient
	targetTable grafanajson.TargetTable
}

// New creates a new Handler
func New(client sciensano.APIClient) (handler *Handler) {
	handler = &Handler{
		Sciensano: client,
	}

	handler.targetTable = grafanajson.TargetTable{
		"hospitalisations":          {TableQueryFunc: handler.buildHospitalisationsResponse},
		"hospitalisations-region":   {TableQueryFunc: handler.buildHospitalisationsResponse},
		"hospitalisations-province": {TableQueryFunc: handler.buildHospitalisationsResponse},
		"hospitalisations-details":  {TableQueryFunc: handler.buildHospitalisationsDetailsResponse},
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

func (handler *Handler) buildHospitalisationsResponse(ctx context.Context, target string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var entries *datasets.Dataset
	switch target {
	case "hospitalisations":
		entries, err = handler.Sciensano.GetHospitalisations(ctx)
	case "hospitalisations-region":
		entries, err = handler.Sciensano.GetHospitalisationsByRegion(ctx)
	case "hospitalisations-province":
		entries, err = handler.Sciensano.GetHospitalisationsByProvince(ctx)
	}

	if err != nil {
		return
	}

	entries.ApplyRange(args.Range.From, args.Range.To)

	timestampColumn := make(grafanajson.TableQueryResponseTimeColumn, 0, len(entries.Timestamps))
	for _, timestamp := range entries.Timestamps {
		timestampColumn = append(timestampColumn, timestamp)
	}

	response = &grafanajson.TableQueryResponse{
		Columns: []grafanajson.TableQueryResponseColumn{{
			Text: "Timestamp",
			Data: timestampColumn,
		}},
	}

	for _, group := range entries.Groups {
		dataColumn := make(grafanajson.TableQueryResponseNumberColumn, 0, len(group.Values))
		for _, value := range group.Values {
			dataColumn = append(dataColumn, float64(value.(*sciensano.HospitalisationsEntry).In))
		}

		name := group.Name
		if name == "" {
			if target == "hospitalisation" {
				name = "hospitalisation"
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

func (handler *Handler) buildHospitalisationsDetailsResponse(ctx context.Context, _ string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var entries *datasets.Dataset
	entries, err = handler.Sciensano.GetHospitalisations(ctx)

	if err != nil {
		return
	}

	entries.ApplyRange(args.Range.From, args.Range.To)

	timestampColumn := make(grafanajson.TableQueryResponseTimeColumn, 0, len(entries.Timestamps))
	for _, timestamp := range entries.Timestamps {
		timestampColumn = append(timestampColumn, timestamp)
	}

	dataColumns := []grafanajson.TableQueryResponseNumberColumn{
		make(grafanajson.TableQueryResponseNumberColumn, len(timestampColumn)),
		make(grafanajson.TableQueryResponseNumberColumn, len(timestampColumn)),
		make(grafanajson.TableQueryResponseNumberColumn, len(timestampColumn)),
		make(grafanajson.TableQueryResponseNumberColumn, len(timestampColumn)),
	}

	for index, value := range entries.Groups[0].Values {
		dataColumns[0][index] = float64(value.(*sciensano.HospitalisationsEntry).In)
		dataColumns[1][index] = float64(value.(*sciensano.HospitalisationsEntry).InICU)
		dataColumns[2][index] = float64(value.(*sciensano.HospitalisationsEntry).InResp)
		dataColumns[3][index] = float64(value.(*sciensano.HospitalisationsEntry).InECMO)
	}

	response = &grafanajson.TableQueryResponse{
		Columns: []grafanajson.TableQueryResponseColumn{
			{
				Text: "Timestamp",
				Data: timestampColumn,
			},
			{
				Text: "in",
				Data: dataColumns[0],
			},
			{
				Text: "inICU",
				Data: dataColumns[1],
			},
			{
				Text: "inRESP",
				Data: dataColumns[2],
			},
			{
				Text: "inECMO",
				Data: dataColumns[3],
			},
		},
	}
	return
}
