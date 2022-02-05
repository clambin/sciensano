package responder

import (
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/simplejson/v3/query"
)

// GenerateTableQueryResponse generates a TableQueryResponse from a Dataset
func GenerateTableQueryResponse(input *datasets.Dataset, args query.Args) (response *query.TableResponse) {
	input.FilterByRange(args.Range.From, args.Range.To)

	response = &query.TableResponse{
		Columns: []query.Column{{
			Text: "timestamp",
			Data: query.TimeColumn(input.GetTimestamps()),
		}},
	}

	for _, column := range input.GetColumns() {
		values, _ := input.GetValues(column)
		if column == "" {
			column = "(unknown)"
		}
		response.Columns = append(response.Columns, query.Column{
			Text: column,
			Data: query.NumberColumn(values),
		})
	}
	return
}
