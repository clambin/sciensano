package responder

import (
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/simplejson"
)

// GenerateTableQueryResponse generates a TableQueryResponse from a Dataset
func GenerateTableQueryResponse(input *datasets.Dataset, args *simplejson.TableQueryArgs) (response *simplejson.TableQueryResponse) {
	input.ApplyRange(args.Range.From, args.Range.To)

	timestampColumn := make(simplejson.TableQueryResponseTimeColumn, len(input.Timestamps))
	for index, timestamp := range input.Timestamps {
		timestampColumn[index] = timestamp
	}

	response = &simplejson.TableQueryResponse{
		Columns: []simplejson.TableQueryResponseColumn{{
			Text: "timestamp",
			Data: timestampColumn,
		}},
	}

	for _, group := range input.Groups {
		dataColumn := make(simplejson.TableQueryResponseNumberColumn, len(group.Values))
		for index, value := range group.Values {
			dataColumn[index] = value
		}

		name := group.Name
		if name == "" {
			name = "(unknown)"
		}

		response.Columns = append(response.Columns, simplejson.TableQueryResponseColumn{
			Text: name,
			Data: dataColumn,
		})
	}
	return
}
