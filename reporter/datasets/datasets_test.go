package datasets_test

import (
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAltDataset_Copy(t *testing.T) {
	ds := &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 6, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "A", Values: []float64{1, 2, 3, 4, 5, 6}},
			{Name: "B", Values: []float64{11, 21, 31, 41, 51, 61}},
		},
	}

	ds2 := ds.Copy()

	assert.False(t, ds2 == ds)
	assert.Equal(t, ds2, ds)
}

func TestDataset_ApplyRange(t *testing.T) {
	ds := &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 6, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "A", Values: []float64{1, 2, 3, 4, 5, 6}},
			{Name: "B", Values: []float64{11, 21, 31, 41, 51, 61}},
		},
	}

	ds.ApplyRange(time.Time{}, time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC))

	assert.Len(t, ds.Timestamps, 5)
	require.Len(t, ds.Groups, 2)
	assert.Equal(t, "A", ds.Groups[0].Name)
	assert.Equal(t, []float64{1, 2, 3, 4, 5}, ds.Groups[0].Values)
	assert.Equal(t, "B", ds.Groups[1].Name)
	assert.Equal(t, []float64{11, 21, 31, 41, 51}, ds.Groups[1].Values)

	ds.ApplyRange(time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC), time.Time{})

	assert.Len(t, ds.Timestamps, 4)
	require.Len(t, ds.Groups, 2)
	assert.Equal(t, "A", ds.Groups[0].Name)
	assert.Equal(t, []float64{2, 3, 4, 5}, ds.Groups[0].Values)
	assert.Equal(t, "B", ds.Groups[1].Name)
	assert.Equal(t, []float64{21, 31, 41, 51}, ds.Groups[1].Values)
}

func TestAltDataset_Accumulate(t *testing.T) {
	ds := &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 6, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "A", Values: []float64{1, 2, 3, 4, 5, 6}},
			{Name: "B", Values: []float64{11, 21, 31, 41, 51, 61}},
		},
	}

	ds.Accumulate()

	assert.Equal(t, datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 6, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "A", Values: []float64{1, 3, 6, 10, 15, 21}},
			{Name: "B", Values: []float64{11, 32, 63, 104, 155, 216}},
		},
	}, *ds)
}
