package datasets_test

import (
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type groupable struct {
	timestamp time.Time
	group     string
	attrib1   float64
	attrib2   float64
}

func (g groupable) GetTimestamp() time.Time {
	return g.timestamp
}

func (g groupable) GetGroupFieldValue(_ int) string {
	return g.group
}

func (g groupable) GetTotalValue() float64 {
	return g.attrib1 + g.attrib2
}

func (g groupable) GetAttributeNames() []string {
	return []string{"Attrib1", "Attrib2"}
}

func (g groupable) GetAttributeValues() []float64 {
	return []float64{g.attrib1, g.attrib2}
}

func TestGroupMeasurements(t *testing.T) {
	m := []measurement.Measurement{
		&groupable{timestamp: time.Time{}, group: "A", attrib1: 10, attrib2: 20},
		&groupable{timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), group: "A", attrib1: 1, attrib2: 0},
		&groupable{timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), group: "B", attrib1: 0, attrib2: 1},
		&groupable{timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), group: "B", attrib1: 1, attrib2: 0},
		&groupable{timestamp: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), group: "C", attrib1: 0, attrib2: 3},
	}
	ds := datasets.GroupMeasurements(m)

	require.Len(t, ds.Timestamps, 2)
	require.Len(t, ds.Groups, 2)
	assert.Equal(t, datasets.DatasetGroup{Name: "Attrib1", Values: []float64{2, 0}}, ds.Groups[0])
	assert.Equal(t, datasets.DatasetGroup{Name: "Attrib2", Values: []float64{1, 3}}, ds.Groups[1])
}

func TestGroupMeasurementsByType(t *testing.T) {
	m := []measurement.Measurement{
		&groupable{timestamp: time.Time{}, group: "A", attrib1: 10, attrib2: 20},
		&groupable{timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), group: "A", attrib1: 1, attrib2: 0},
		&groupable{timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), group: "B", attrib1: 0, attrib2: 1},
		&groupable{timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), group: "B", attrib1: 1, attrib2: 0},
		&groupable{timestamp: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), group: "C", attrib1: 0, attrib2: 3},
	}
	ds := datasets.GroupMeasurementsByType(m, 1)

	require.Len(t, ds.Timestamps, 2)
	require.Len(t, ds.Groups, 3)
	assert.Equal(t, datasets.DatasetGroup{Name: "A", Values: []float64{1, 0}}, ds.Groups[0])
	assert.Equal(t, datasets.DatasetGroup{Name: "B", Values: []float64{2, 0}}, ds.Groups[1])
	assert.Equal(t, datasets.DatasetGroup{Name: "C", Values: []float64{0, 3}}, ds.Groups[2])
}

func TestGroupMeasurements_Vaccinations(t *testing.T) {
	m := []measurement.Measurement{
		&sciensano.APIVaccinationsResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "85+",
			Dose:      "A",
			Count:     10,
		},
	}
	ds := datasets.GroupMeasurements(m)
	require.Len(t, ds.Timestamps, 1)
}

func bigResponse() []measurement.Measurement {
	input := make([]measurement.Measurement, 0)
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Flanders", "Wallonia", "Brussels"} {
			for _, ageGroup := range []string{"00-11", "12-17", "18-34", "35-44", "45-54", "55-64", "65-74", "75-84", "85+"} {
				for _, dose := range []string{"A", "B", "C", "E"} {
					input = append(input, &sciensano.APIVaccinationsResponseEntry{
						TimeStamp: sciensano.TimeStamp{Time: timestamp},
						Region:    region,
						AgeGroup:  ageGroup,
						Dose:      dose,
						Count:     i,
					})
				}
			}
		}
		timestamp = timestamp.Add(24 * time.Hour)
	}
	return input
}

func BenchmarkGroupMeasurements(b *testing.B) {
	input := bigResponse()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = datasets.GroupMeasurements(input)
	}
}

func BenchmarkGroupMeasurementsByType(b *testing.B) {
	input := bigResponse()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = datasets.GroupMeasurementsByType(input, measurement.GroupByAgeGroup)
	}
}
