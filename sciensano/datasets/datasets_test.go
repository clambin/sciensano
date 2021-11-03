package datasets_test

import (
	"github.com/clambin/sciensano/sciensano/datasets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type testData struct {
	Value int
}

func (ts *testData) Copy() datasets.Copyable {
	return &testData{Value: ts.Value}
}

func TestDataset_Copy(t *testing.T) {
	input := datasets.Dataset{
		Timestamps: []time.Time{time.Now()},
		Groups: []datasets.GroupedDatasetEntry{
			{
				Name:   "test",
				Values: []datasets.Copyable{&testData{Value: 1}},
			},
		},
	}
	output := input.Copy()

	require.NotNil(t, output)
	require.Len(t, output.Timestamps, 1)
	require.Len(t, output.Groups, 1)
	require.Len(t, output.Groups[0].Values, 1)
	assert.Equal(t, 1, output.Groups[0].Values[0].(*testData).Value)
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
		Groups: []datasets.GroupedDatasetEntry{
			{
				Name: "A",
				Values: []datasets.Copyable{
					&testData{Value: 1},
					&testData{Value: 2},
					&testData{Value: 3},
					&testData{Value: 4},
					&testData{Value: 5},
					&testData{Value: 6},
				},
			},
			{
				Name: "B",
				Values: []datasets.Copyable{
					&testData{Value: 11},
					&testData{Value: 21},
					&testData{Value: 31},
					&testData{Value: 41},
					&testData{Value: 51},
					&testData{Value: 61},
				},
			},
		},
	}

	ds.ApplyRange(time.Time{}, time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC))

	assert.Len(t, ds.Timestamps, 5)
	require.Len(t, ds.Groups, 2)
	assert.Equal(t, "A", ds.Groups[0].Name)
	require.Len(t, ds.Groups[0].Values, 5)
	assert.Equal(t, 1, ds.Groups[0].Values[0].(*testData).Value)
	assert.Equal(t, 5, ds.Groups[0].Values[4].(*testData).Value)
	assert.Equal(t, "B", ds.Groups[1].Name)
	require.Len(t, ds.Groups[1].Values, 5)
	assert.Equal(t, 11, ds.Groups[1].Values[0].(*testData).Value)
	assert.Equal(t, 51, ds.Groups[1].Values[4].(*testData).Value)

	ds.ApplyRange(time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC), time.Time{})

	assert.Len(t, ds.Timestamps, 4)
	require.Len(t, ds.Groups, 2)
	assert.Equal(t, "A", ds.Groups[0].Name)
	require.Len(t, ds.Groups[0].Values, 4)
	assert.Equal(t, 2, ds.Groups[0].Values[0].(*testData).Value)
	assert.Equal(t, 5, ds.Groups[0].Values[3].(*testData).Value)
	assert.Equal(t, "B", ds.Groups[1].Name)
	require.Len(t, ds.Groups[1].Values, 4)
	assert.Equal(t, 21, ds.Groups[1].Values[0].(*testData).Value)
	assert.Equal(t, 51, ds.Groups[1].Values[3].(*testData).Value)
}
