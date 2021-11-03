package datasets

import "time"

// Dataset is the response of the GetVaccinations
type Dataset struct {
	Timestamps []time.Time
	Groups     []GroupedDatasetEntry
}

// GroupedDatasetEntry contains the values for the (grouped) vaccinations
type GroupedDatasetEntry struct {
	Name   string
	Values []Copyable
}

// Copyable interface for GroupedDatasetEntry structure provides a copy of a grouped value
type Copyable interface {
	Copy() Copyable
}

// Copy makes a copy of a Dataset
func (dataset Dataset) Copy() (result *Dataset) {
	result = &Dataset{
		Timestamps: make([]time.Time, len(dataset.Timestamps)),
	}
	for index, timestamp := range dataset.Timestamps {
		result.Timestamps[index] = timestamp
	}
	result.Groups = make([]GroupedDatasetEntry, len(dataset.Groups))
	for index, group := range dataset.Groups {
		entry := GroupedDatasetEntry{
			Name:   group.Name,
			Values: make([]Copyable, len(group.Values)),
		}
		for index2, value := range group.Values {
			entry.Values[index2] = value.Copy()
		}
		result.Groups[index] = entry
	}
	return
}

// ApplyRange limits the data in a dataset to the provided from/to timestamps. If to/from is zero, it will be ignored.
func (dataset *Dataset) ApplyRange(from, to time.Time) {
	first := 0
	if from.IsZero() == false {
		for index, timestamp := range dataset.Timestamps {
			if !timestamp.Before(from) {
				first = index
				break
			}
		}
	}
	if first != 0 {
		dataset.Timestamps = dataset.Timestamps[first:]
		for index, group := range dataset.Groups {
			dataset.Groups[index].Values = group.Values[first:]
		}
	}

	last := len(dataset.Timestamps)
	if to.IsZero() == false {
		for index, timestamp := range dataset.Timestamps {
			if timestamp.After(to) {
				break
			}
			last = index
		}
	}
	if last != len(dataset.Timestamps) {
		dataset.Timestamps = dataset.Timestamps[:last+1]
		for index, group := range dataset.Groups {
			dataset.Groups[index].Values = group.Values[:last+1]
		}
	}
}
