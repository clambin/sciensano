package hospitalisations

import (
	"context"
	"fmt"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler/response"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	log "github.com/sirupsen/logrus"
	"time"
)

// Handler implements a grafana-json handler for COVID-19 cases
type Handler struct {
	Sciensano   reporter.Reporter
	targetTable grafanajson.TargetTable
}

// New creates a new Handler
func New(client reporter.Reporter) (handler *Handler) {
	handler = &Handler{
		Sciensano: client,
	}

	handler.targetTable = grafanajson.TargetTable{
		"hospitalisations":          {TableQueryFunc: handler.buildHospitalisationsResponse},
		"hospitalisations-region":   {TableQueryFunc: handler.buildHospitalisationsResponse},
		"hospitalisations-province": {TableQueryFunc: handler.buildHospitalisationsResponse},
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
	log.WithFields(log.Fields{"duration": time.Now().Sub(start), "target": target}).Info("TableQuery called")
	return
}

func (handler *Handler) buildHospitalisationsResponse(ctx context.Context, target string, args *grafanajson.TableQueryArgs) (output *grafanajson.TableQueryResponse, err error) {
	var entries *datasets.Dataset
	switch target {
	case "hospitalisations":
		entries, err = handler.Sciensano.GetHospitalisations(ctx)
	case "hospitalisations-region":
		entries, err = handler.Sciensano.GetHospitalisationsByRegion(ctx)
	case "hospitalisations-province":
		entries, err = handler.Sciensano.GetHospitalisationsByProvince(ctx)
	}

	if err == nil {
		output = response.GenerateTableQueryResponse(entries, args)
		return
	}

	return
}
