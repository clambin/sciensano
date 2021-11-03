package sciensano

import "github.com/clambin/sciensano/sciensano/datasets"

// CasesEntry contains the cases for a single timestamp
type CasesEntry struct {
	Count int
}

// Copy makes a copy of a CasesEntry
func (entry *CasesEntry) Copy() datasets.Copyable {
	return &CasesEntry{Count: entry.Count}
}
