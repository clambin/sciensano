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
func (entry *HospitalisationsEntry) Add(input apiclient.APIHospitalisationsResponseEntry) {
	entry.In += input.TotalIn
	entry.InICU += input.TotalInICU
	entry.InResp += input.TotalInResp
	entry.InECMO += input.TotalInECMO
}
