package sciensano

import (
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/datasets"
)

// HospitalisationsEntry contains the hospitalisation values for a single timestamp
type HospitalisationsEntry struct {
	In     int
	InICU  int
	InResp int
	InECMO int
}

// NewHospitalisationsEntry returns a new HospitalisationsEntry, as a GroupedEntry. Used by groupMeasurements
func NewHospitalisationsEntry() GroupedEntry {
	return &HospitalisationsEntry{}
}

// Copy makes a copy of a HospitalisationsEntry
func (entry *HospitalisationsEntry) Copy() datasets.Copyable {
	return &HospitalisationsEntry{
		In:     entry.In,
		InICU:  entry.InICU,
		InResp: entry.InResp,
		InECMO: entry.InECMO,
	}
}

// Add adds the passed HospitalisationEntry values to its own values
func (entry *HospitalisationsEntry) Add(input apiclient.Measurement) {
	entry.In += input.(*apiclient.APIHospitalisationsResponseEntry).TotalIn
	entry.InICU += input.(*apiclient.APIHospitalisationsResponseEntry).TotalInICU
	entry.InResp += input.(*apiclient.APIHospitalisationsResponseEntry).TotalInResp
	entry.InECMO += input.(*apiclient.APIHospitalisationsResponseEntry).TotalInECMO
}
