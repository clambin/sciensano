package sciensano

import (
	"github.com/clambin/sciensano/sciensano/datasets"
)

// VaccinationsEntry contains the vaccination values for a single timestamp
type VaccinationsEntry struct {
	Partial    int
	Full       int
	SingleDose int
	Booster    int
}

// Copy makes a copy of a VaccinationsEntry
func (entry *VaccinationsEntry) Copy() datasets.Copyable {
	return &VaccinationsEntry{
		Partial:    entry.Partial,
		Full:       entry.Full,
		SingleDose: entry.SingleDose,
		Booster:    entry.Booster,
	}
}

// Total calculates the total number of vaccinations for one entry
func (entry VaccinationsEntry) Total() int {
	return entry.Partial + entry.Full + entry.SingleDose + entry.Booster
}

const (
	VaccinationTypePartial int = iota
	VaccinationTypeFull
	VaccinationTypeBooster
)

// GetValue returns the vaccination count for the specified type
func (entry VaccinationsEntry) GetValue(vaccinationType int) (value int) {
	switch vaccinationType {
	case VaccinationTypePartial:
		value = entry.Partial
	case VaccinationTypeFull:
		value = entry.Full + entry.SingleDose
	case VaccinationTypeBooster:
		value = entry.Booster
	}

	return
}
