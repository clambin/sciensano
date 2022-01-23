package vaccinations

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson"
)

type RateHandler struct {
	reporter.Reporter
	reporter.VaccinationType
	Scope
	demographics.Demographics
	helper *GroupedHandler
}

var _ simplejson.Handler = &RateHandler{}

func (r RateHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{TableQuery: r.tableQuery}
}

func (r *RateHandler) tableQuery(ctx context.Context, args *simplejson.TableQueryArgs) (response *simplejson.TableQueryResponse, err error) {
	if r.helper == nil {
		r.helper = &GroupedHandler{
			Reporter:        r.Reporter,
			VaccinationType: r.VaccinationType,
			Scope:           r.Scope,
		}
	}

	response, err = r.helper.tableQuery(ctx, args)

	if err != nil {
		return nil, fmt.Errorf("vaccination rate failed: %w", err)
	}

	response.Columns = filterUnknownColumns(response.Columns)

	switch r.Scope {
	case ScopeAge:
		ageGroupFigures := r.Demographics.GetAgeGroupFigures()
		prorateFigures(response, ageGroupFigures)
	case ScopeRegion:
		regionFigures := r.Demographics.GetRegionFigures()
		// demographics counts figures for Ostbelgien as part of Wallonia. Hardcode the split here.
		// yes, it's ugly. :-)
		_, ok := regionFigures["Ostbelgien"]
		if !ok {
			population, _ := regionFigures["Wallonia"]
			population -= 78000
			regionFigures["Wallonia"] = population
			regionFigures["Ostbelgien"] = 78000
		}
		prorateFigures(response, regionFigures)
	}
	return
}

func filterUnknownColumns(columns []simplejson.TableQueryResponseColumn) []simplejson.TableQueryResponseColumn {
	newColumns := make([]simplejson.TableQueryResponseColumn, 0, len(columns))
	shouldReplace := false
	for _, column := range columns {
		if column.Text == "(unknown)" {
			shouldReplace = true
			continue
		}
		newColumns = append(newColumns, column)
	}
	if shouldReplace {
		return newColumns
	}
	return columns
}

func prorateFigures(result *simplejson.TableQueryResponse, groups map[string]int) {
	for _, column := range result.Columns {
		switch data := column.Data.(type) {
		case simplejson.TableQueryResponseNumberColumn:
			figure, ok := groups[column.Text]
			for index, entry := range data {
				if ok && figure != 0 {
					data[index] = entry / float64(figure)
				} else {
					data[index] = 0
				}
			}
		}
	}
}
