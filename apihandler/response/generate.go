package response

import (
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/simplejson"
)

// GenerateTableQueryResponse generates a TableQueryResponse from a Dataset
func GenerateTableQueryResponse(input *datasets.Dataset, args *simplejson.TableQueryArgs) (response *simplejson.TableQueryResponse) {
	input.ApplyRange(args.Range.From, args.Range.To)

	timestampColumn := make(simplejson.TableQueryResponseTimeColumn, 0, len(input.Timestamps))
	for _, timestamp := range input.Timestamps {
		timestampColumn = append(timestampColumn, timestamp)
	}

	response = &simplejson.TableQueryResponse{
		Columns: []simplejson.TableQueryResponseColumn{{
			Text: "timestamp",
			Data: timestampColumn,
		}},
	}

	for _, group := range input.Groups {
		dataColumn := make(simplejson.TableQueryResponseNumberColumn, 0, len(group.Values))
		for _, value := range group.Values {
			dataColumn = append(dataColumn, value)
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
