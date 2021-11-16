package reporter

import (
	"fmt"
	"github.com/clambin/sciensano/reporter/datasets"
)

// VaccinesGetter contains all required methods to retrieve vaccines data
type VaccinesGetter interface {
	GetVaccines() (results *datasets.Dataset, err error)
}

// GetVaccines returns all vaccines data
func (client *Client) GetVaccines() (results *datasets.Dataset, err error) {
	return client.Cache.MaybeGenerate("Vaccines", func() (output *datasets.Dataset, err2 error) {
		if apiResult, ok := client.Vaccines.Get("Vaccines"); ok {
			output = datasets.GroupMeasurements(apiResult)
		} else {
			err2 = fmt.Errorf("cache does not contain Vaccines entries")
		}
		return
	})
}
