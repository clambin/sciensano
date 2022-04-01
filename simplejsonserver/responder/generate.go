package responder

import (
	"github.com/clambin/simplejson/v3/dataset"
	"github.com/clambin/simplejson/v3/query"
)

// GenerateTableQueryResponse generates a TableQueryResponse from a Dataset
func GenerateTableQueryResponse(input *dataset.Dataset, args query.Args) (response *query.TableResponse) {
	input.FilterByRange(args.Range.From, args.Range.To)
	return input.GenerateTableResponse()
}
