package datasets

import "time"

// Dataset is the processed responder of a sciensano/vaccine API call
type Dataset struct {
	Timestamps []time.Time
	Groups     []DatasetGroup
}

// DatasetGroup is one group of values of an Dataset
type DatasetGroup struct {
	Name   string
	Values []float64
}

// Copy returns a copy of a data set
func (dataset *Dataset) Copy() (result *Dataset) {
	result = &Dataset{
		Timestamps: make([]time.Time, len(dataset.Timestamps)),
		Groups:     make([]DatasetGroup, len(dataset.Groups)),
	}

	copy(result.Timestamps, dataset.Timestamps)
	for index, group := range dataset.Groups {
		result.Groups[index].Name = group.Name
		result.Groups[index].Values = make([]float64, len(group.Values))
		copy(result.Groups[index].Values, group.Values)
	}
	return
}

// ApplyRange limits the data in a dataset to the provided from/to timestamps. If to/from is zero, it will be ignored.
func (dataset *Dataset) ApplyRange(from, to time.Time) {
	if from.IsZero() == false {
		first := 0
		for index, timestamp := range dataset.Timestamps {
			if !timestamp.Before(from) {
				first = index
				break
			}
		}
		if first != 0 {
			dataset.Timestamps = dataset.Timestamps[first:]
			for index, group := range dataset.Groups {
				dataset.Groups[index].Values = group.Values[first:]
			}
		}
	}

	if to.IsZero() == false {
		last := len(dataset.Timestamps)
		for index, timestamp := range dataset.Timestamps {
			if timestamp.After(to) {
				break
			}
			last = index
		}
		if last != len(dataset.Timestamps) {
			dataset.Timestamps = dataset.Timestamps[:last+1]
			for index, group := range dataset.Groups {
				dataset.Groups[index].Values = group.Values[:last+1]
			}
		}
	}
}

// Accumulate accumulates the values in each of the dataset's groups
func (dataset *Dataset) Accumulate() {
	for _, group := range dataset.Groups {
		var accu float64
		for index, entry := range group.Values {
			accu += entry
			group.Values[index] = accu
		}
	}
}
