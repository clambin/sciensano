package vaccinations

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/demographics/bracket"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/query"
)

type RateHandler struct {
	reporter.Reporter
	reporter.VaccinationType
	Scope
	demographics.Fetcher
	helper *GroupedHandler
}

var _ simplejson.Handler = &RateHandler{}

func (r RateHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: r.tableQuery}
}

func (r *RateHandler) tableQuery(ctx context.Context, req query.Request) (response query.Response, err error) {
	if r.helper == nil {
		r.helper = &GroupedHandler{
			Reporter:        r.Reporter,
			VaccinationType: r.VaccinationType,
			Scope:           r.Scope,
		}
	}

	response, err = r.helper.tableQuery(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("vaccination rate failed: %w", err)
	}

	resp := response.(*query.TableResponse)
	resp.Columns = filterUnknownColumns(resp.Columns)

	var figures map[string]int
	switch r.Scope {
	case ScopeRegion:
		figures = r.Fetcher.GetByRegion()
	case ScopeAge:
		figures = make(map[string]int)
		for _, column := range resp.Columns {
			if column.Text == "timestamp" {
				continue
			}
			var b bracket.Bracket
			b, err = bracket.FromString(column.Text)
			if err != nil {
				return nil, fmt.Errorf("invalid age bracket: '%s' : %w", column.Text, err)
			}
			figures[column.Text] = r.Fetcher.GetByAgeBracket(b)
		}
	}
	if err == nil {
		prorateFigures(resp, figures)
	}
	return
}

func filterUnknownColumns(columns []query.Column) []query.Column {
	newColumns := make([]query.Column, 0, len(columns))
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

func prorateFigures(result *query.TableResponse, groups map[string]int) {
	for _, column := range result.Columns {
		switch data := column.Data.(type) {
		case query.NumberColumn:
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
