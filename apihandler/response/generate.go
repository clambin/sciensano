package response

import (
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/reporter/datasets"
)

// GenerateTableQueryResponse generates a TableQueryResponse from a Dataset
func GenerateTableQueryResponse(input *datasets.Dataset, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse) {
	input.ApplyRange(args.Range.From, args.Range.To)

	timestampColumn := make(grafanajson.TableQueryResponseTimeColumn, 0, len(input.Timestamps))
	for _, timestamp := range input.Timestamps {
		timestampColumn = append(timestampColumn, timestamp)
	}

	response = &grafanajson.TableQueryResponse{
		Columns: []grafanajson.TableQueryResponseColumn{{
			Text: "timestamp",
			Data: timestampColumn,
		}},
	}

	for _, group := range input.Groups {
		dataColumn := make(grafanajson.TableQueryResponseNumberColumn, 0, len(group.Values))
		for _, value := range group.Values {
			dataColumn = append(dataColumn, value)
		}

		name := group.Name
		if name == "" {
			name = "(unknown)"
		}

		response.Columns = append(response.Columns, grafanajson.TableQueryResponseColumn{
			Text: name,
			Data: dataColumn,
		})
	}
	return
}
