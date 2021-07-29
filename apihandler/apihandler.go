package apihandler

import (
	"fmt"
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/vaccines"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Handler implements a Grafana SimpleJson API Handler that gets BE covid stats
type Handler struct {
	Sciensano   sciensano.API
	Vaccines    *vaccines.Server
	lastDate    map[string]time.Time
	targetTable TargetTable
}

// Create a Handler
func Create() (*Handler, error) {
	handler := Handler{
		Sciensano: &sciensano.Client{
			HTTPClient:    &http.Client{Timeout: 20 * time.Second},
			CacheDuration: 15 * time.Minute,
		},
		Vaccines: vaccines.New(),
	}

	handler.targetTable = TargetTable{
		"tests":                    {tableResponseBuild: handler.buildTestTableResponse},
		"vaccinations":             {tableResponseBuild: handler.buildVaccinationTableResponse},
		"vacc-age-partial":         {tableResponseBuild: handler.buildGroupedVaccinationTableResponse},
		"vacc-age-full":            {tableResponseBuild: handler.buildGroupedVaccinationTableResponse},
		"vacc-age-rate-partial":    {tableResponseBuild: handler.buildGroupedVaccinationRateTableResponse},
		"vacc-age-rate-full":       {tableResponseBuild: handler.buildGroupedVaccinationRateTableResponse},
		"vacc-region-partial":      {tableResponseBuild: handler.buildGroupedVaccinationTableResponse},
		"vacc-region-full":         {tableResponseBuild: handler.buildGroupedVaccinationTableResponse},
		"vacc-region-rate-partial": {tableResponseBuild: handler.buildGroupedVaccinationRateTableResponse},
		"vacc-region-rate-full":    {tableResponseBuild: handler.buildGroupedVaccinationRateTableResponse},
		"vaccination-lag":          {tableResponseBuild: handler.buildVaccinationLagTableResponse},
		"vaccines":                 {tableResponseBuild: handler.buildVaccineTableResponse},
		"vaccines-stats":           {tableResponseBuild: handler.buildVaccineStatsTableResponse},
		"vaccines-time":            {tableResponseBuild: handler.buildVaccineTimeTableResponse},
	}

	return &handler, nil
}

// Endpoints tells the server which endpoints we have implemented
func (handler *Handler) Endpoints() grafana_json.Endpoints {
	return grafana_json.Endpoints{
		Search:      handler.Search,
		TableQuery:  handler.TableQuery,
		Annotations: handler.Annotations,
	}
}

type SeriesResponseBuildFunc func(begin, end time.Time, target string) (response *grafana_json.QueryResponse)
type TableResponseBuildFunc func(begin, end time.Time, target string) (response *grafana_json.TableQueryResponse)

type TargetTable map[string]struct {
	seriesResponseBuild SeriesResponseBuildFunc
	tableResponseBuild  TableResponseBuildFunc
}

// Search returns all supported targets
func (handler *Handler) Search() (targets []string) {
	for target := range handler.targetTable {
		targets = append(targets, target)
	}
	sort.Strings(targets)
	return
}

func (handler *Handler) TableQuery(target string, args *grafana_json.TableQueryArgs) (response *grafana_json.TableQueryResponse, err error) {

	builder, ok := handler.targetTable[target]

	if ok == false || builder.tableResponseBuild == nil {
		return nil, fmt.Errorf("unknown target '%s'", target)
	}

	response = builder.tableResponseBuild(args.Range.From, args.Range.To, target)

	// log if there's a new update
	if response != nil {
		handler.logUpdates(target, response)
	}

	return
}

func (handler *Handler) logUpdates(target string, response *grafana_json.TableQueryResponse) {
	if handler.lastDate == nil {
		handler.lastDate = make(map[string]time.Time)
	}

	key := target
	if strings.HasPrefix(target, "vaccines") {
		key = "vaccines"
	} else if strings.HasPrefix(target, "vacc") {
		key = "vaccinations"
	}

	latest := time.Time{}
	for _, column := range response.Columns {
		if column.Text == "timestamp" {
			timestamps := column.Data.(grafana_json.TableQueryResponseTimeColumn)
			if len(timestamps) > 0 {
				latest = timestamps[len(timestamps)-1]
				break
			}
		}
	}

	if entry, ok := handler.lastDate[key]; ok == false {
		handler.lastDate[key] = latest
	} else if latest.After(entry) {
		handler.lastDate[key] = latest
		log.WithFields(log.Fields{"target": key, "time": latest}).Info("new data found")
	}
}
