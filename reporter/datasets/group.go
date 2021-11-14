package datasets

import (
	"github.com/clambin/sciensano/measurement"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

// GroupMeasurements groups a slice of measurements into a dataset.  The measurement's attributes (as determined by
// the measurements' GetAttributeNames function) will be stored in the dataset's columns.
//
// e.g. if the measurement's GetAttributeNames returns []string{"A", "B"}, the dataset will look as follows:
//		&Dataset{
//			Timestamps: []time.Time{...},
//			Columns: []DatasetGroup{
//			{
//				Name: "A", Values: []float64{...},
//				Name: "B", Values: []float64{...},
//			},
//		}
func GroupMeasurements(entries []measurement.Measurement) (results *Dataset) {
	before := time.Now()
	defer log.WithField("time", time.Now().Sub(before)).Debug("GroupMeasurementDetails")

	results = &Dataset{
		Timestamps: make([]time.Time, 0),
	}

	currentTimestamp := time.Time{}
	currentIndex := -1

	for _, entry := range entries {
		if results.Groups == nil {
			groups := entry.GetAttributeNames()
			results.Groups = make([]DatasetGroup, len(groups))
			for index, group := range groups {
				results.Groups[index] = DatasetGroup{
					Name:   group,
					Values: make([]float64, 0),
				}
			}
		}

		if entry.GetTimestamp().IsZero() {
			continue
		}

		if currentTimestamp.IsZero() || !currentTimestamp.Equal(entry.GetTimestamp()) {
			results.Timestamps = append(results.Timestamps, entry.GetTimestamp())
			for index := range entry.GetAttributeNames() {
				results.Groups[index].Values = append(results.Groups[index].Values, 0)
			}
			currentIndex++
			currentTimestamp = entry.GetTimestamp()
		}

		for index, value := range entry.GetAttributeValues() {
			results.Groups[index].Values[currentIndex] += value
		}
	}

	return
}

// GroupMeasurementsByType groups a slice of measurements into a dataset by the measurement's groupField
//
// e.g. if the values for the groupField are "A" and "B", the dataset will look as follows:
//		&Dataset{
//			Timestamps: []time.Time{...},
//			Columns: []DatasetGroup{
//			{
//				Name: "A", Values: []float64{...},
//				Name: "B", Values: []float64{...},
//			},
//		}
func GroupMeasurementsByType(entries []measurement.Measurement, groupField int) (results *Dataset) {
	groupFieldValues, groupFieldValuesIndex := getUniqueSortedGroupValues(entries, groupField)

	results = &Dataset{
		Timestamps: make([]time.Time, 0),
		Groups:     make([]DatasetGroup, len(groupFieldValues)),
	}

	for index, groupValue := range groupFieldValues {
		results.Groups[index] = DatasetGroup{
			Name:   groupValue,
			Values: make([]float64, 0),
		}
	}

	currentTimestamp := time.Time{}
	count := 0
	for _, entry := range entries {
		if entry.GetTimestamp().IsZero() {
			continue
		}

		if currentTimestamp.IsZero() || entry.GetTimestamp().Equal(currentTimestamp) == false {
			currentTimestamp = entry.GetTimestamp()
			results.Timestamps = append(results.Timestamps, currentTimestamp)
			count++
			for index := range results.Groups {
				results.Groups[index].Values = append(results.Groups[index].Values, 0.0)
			}
		}

		results.Groups[groupFieldValuesIndex[entry.GetGroupFieldValue(groupField)]].Values[count-1] += entry.GetTotalValue()
	}

	return
}

// getUniqueSortedGroupValues returns the sorted unique values of the measurement's groupField
func getUniqueSortedGroupValues(input []measurement.Measurement, groupField int) (groupNames []string, groupNameIndex map[string]int) {
	groupNameIndex = make(map[string]int)
	for _, entry := range input {
		groupNameIndex[entry.GetGroupFieldValue(groupField)] = 0
	}
	groupNames = make([]string, 0, len(groupNameIndex))
	for key := range groupNameIndex {
		groupNames = append(groupNames, key)
	}
	sort.Strings(groupNames)
	for index, value := range groupNames {
		groupNameIndex[value] = index
	}
	return
}
