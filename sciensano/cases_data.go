package sciensano

import (
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/datasets"
)

// CasesEntry contains the cases for a single timestamp
type CasesEntry struct {
	Count int
}

func NewCasesEntry() GroupedEntry {
	return &CasesEntry{}
}

// Copy makes a copy of a CasesEntry
func (entry *CasesEntry) Copy() datasets.Copyable {
	return &CasesEntry{Count: entry.Count}
}

// Add adds the passed CasesEntry values to its own values
func (entry *CasesEntry) Add(input apiclient.Measurement) {
	entry.Count += input.(*apiclient.APICasesResponseEntry).Cases
}
