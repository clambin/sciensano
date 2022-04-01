package reporter

import (
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/simplejson/v3/dataset"
)

func NewFromAPIResponse(response []apiclient.APIResponse) (d *dataset.Dataset) {
	d = dataset.New()
	if len(response) == 0 {
		return
	}
	// optimization: all responses will be of same type, so will have same attribute names
	attribs := response[0].GetAttributeNames()
	for _, entry := range response {
		ts := entry.GetTimestamp()
		values := entry.GetAttributeValues()

		for index, attrib := range attribs {
			d.Add(ts, attrib, values[index])
		}
	}
	return
}

func NewGroupedFromAPIResponse(response []apiclient.APIResponse, groupField int) (d *dataset.Dataset) {
	d = dataset.New()
	for _, entry := range response {
		d.Add(entry.GetTimestamp(), entry.GetGroupFieldValue(groupField), entry.GetTotalValue())
	}
	return
}
