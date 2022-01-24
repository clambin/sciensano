package responder

import (
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/simplejson/v2/query"
)

// GenerateTableQueryResponse generates a TableQueryResponse from a Dataset
func GenerateTableQueryResponse(input *datasets.Dataset, args query.Args) (response *query.TableResponse) {
	input.ApplyRange(args.Range.From, args.Range.To)

	timestampColumn := make(query.TimeColumn, len(input.Timestamps))
	for index, timestamp := range input.Timestamps {
		timestampColumn[index] = timestamp
	}

	response = &query.TableResponse{
		Columns: []query.Column{{
			Text: "timestamp",
			Data: timestampColumn,
		}},
	}

	for _, group := range input.Groups {
		dataColumn := make(query.NumberColumn, len(group.Values))
		for index, value := range group.Values {
			dataColumn[index] = value
		}

		name := group.Name
		if name == "" {
			name = "(unknown)"
		}

		response.Columns = append(response.Columns, query.Column{
			Text: name,
			Data: dataColumn,
		})
	}
	return
}
