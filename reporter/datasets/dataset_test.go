package datasets_test

import (
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

type testResponse struct {
	timestamp time.Time
	group     string
	value     float64
}

var _ apiclient.APIResponse = &testResponse{}

func (t testResponse) GetTimestamp() time.Time {
	return t.timestamp
}

func (t testResponse) GetGroupFieldValue(groupField int) (value string) {
	if groupField == 0 {
		return ""
	}
	return t.group
}

func (t testResponse) GetTotalValue() float64 {
	return t.value
}

func (t testResponse) GetAttributeNames() []string {
	panic("implement me")
}

func (t testResponse) GetAttributeValues() []float64 {
	panic("implement me")
}

func TestNewGroupedFromAPIResponse(t *testing.T) {
	var input []apiclient.APIResponse
	input = []apiclient.APIResponse{
		&testResponse{
			timestamp: time.Date(2022, 2, 2, 0, 0, 0, 0, time.UTC),
			group:     "A",
			value:     1,
		},
		&testResponse{
			timestamp: time.Date(2022, 2, 2, 0, 0, 0, 0, time.UTC),
			group:     "B",
			value:     2,
		},
		&testResponse{
			timestamp: time.Date(2022, 2, 2, 0, 0, 0, 0, time.UTC),
			group:     "A",
			value:     3,
		},
		&testResponse{
			timestamp: time.Date(2022, 2, 3, 0, 0, 0, 0, time.UTC),
			group:     "A",
			value:     5,
		},
	}

	d := datasets.NewGroupedFromAPIResponse(input, 1)

	assert.Equal(t, []time.Time{
		time.Date(2022, time.February, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.February, 3, 0, 0, 0, 0, time.UTC),
	}, d.GetTimestamps())

	assert.Equal(t, []string{"A", "B"}, d.GetColumns())

	values, ok := d.GetValues("A")
	require.True(t, ok)
	assert.Equal(t, []float64{4, 5}, values)

	values, ok = d.GetValues("B")
	require.True(t, ok)
	assert.Equal(t, []float64{2, 0}, values)

	_, ok = d.GetValues("C")
	assert.False(t, ok)
}

func BenchmarkDataSetNewGroupedFromAPIResponse(b *testing.B) {
	var input []apiclient.APIResponse
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 720; i++ {
		for c := 0; c < 9; c++ {
			for d := 0; d < 12; d++ {
				input = append(input, &testResponse{
					timestamp: timestamp,
					group:     strconv.Itoa(c),
					value:     1,
				})
			}
		}
		timestamp = timestamp.Add(24 * time.Hour)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = datasets.NewGroupedFromAPIResponse(input, 1)
	}
}

type testMultiValueResponse struct {
	timestamp time.Time
	value1    float64
	value2    float64
}

var _ apiclient.APIResponse = &testMultiValueResponse{}

func (t testMultiValueResponse) GetTimestamp() time.Time {
	return t.timestamp
}

func (t testMultiValueResponse) GetGroupFieldValue(_ int) (value string) {
	panic("implement me")
}

func (t testMultiValueResponse) GetTotalValue() float64 {
	panic("implement me")
}

func (t testMultiValueResponse) GetAttributeNames() []string {
	return []string{"Value1", "Value2"}
}

func (t testMultiValueResponse) GetAttributeValues() []float64 {
	return []float64{t.value1, t.value2}
}

func TestNewFromAPIResponse(t *testing.T) {
	var input []apiclient.APIResponse
	input = []apiclient.APIResponse{
		&testMultiValueResponse{
			timestamp: time.Date(2022, 2, 2, 0, 0, 0, 0, time.UTC),
			value1:    1,
			value2:    2,
		},
		&testMultiValueResponse{
			timestamp: time.Date(2022, 2, 2, 0, 0, 0, 0, time.UTC),
			value1:    3,
			value2:    4,
		},
		&testMultiValueResponse{
			timestamp: time.Date(2022, 2, 3, 0, 0, 0, 0, time.UTC),
			value1:    5,
			value2:    6,
		},
	}

	d := datasets.NewFromAPIResponse(input)

	assert.Equal(t, []time.Time{
		time.Date(2022, time.February, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.February, 3, 0, 0, 0, 0, time.UTC),
	}, d.GetTimestamps())

	assert.Equal(t, []string{"Value1", "Value2"}, d.GetColumns())

	values, ok := d.GetValues("Value1")
	require.True(t, ok)
	assert.Equal(t, []float64{4, 5}, values)

	values, ok = d.GetValues("Value2")
	require.True(t, ok)
	assert.Equal(t, []float64{6, 6}, values)
}

func BenchmarkDataSetNewFromAPIResponse(b *testing.B) {
	var input []apiclient.APIResponse
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 720; i++ {
		for c := 0; c < 9; c++ {
			for d := 0; d < 6; d++ {
				input = append(input, &testMultiValueResponse{
					timestamp: timestamp,
					value1:    1,
					value2:    2,
				})
			}
		}
		timestamp = timestamp.Add(24 * time.Hour)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = datasets.NewFromAPIResponse(input)
	}
}

func TestDataset_AddColumn(t *testing.T) {
	d := datasets.New()
	assert.NotNil(t, d)

	for day := 1; day < 5; day++ {
		d.Add(time.Date(2022, time.January, 5-day, 0, 0, 0, 0, time.UTC), "A", float64(5-day))
	}

	d.AddColumn("B", func(values map[string]float64) float64 {
		return values["A"] * 2
	})

	values, ok := d.GetValues("B")
	require.True(t, ok)
	assert.Equal(t, []float64{2, 4, 6, 8}, values)
}

func TestDataset_FilterByRange(t *testing.T) {
	d := datasets.New()
	assert.NotNil(t, d)

	for day := 1; day < 32; day++ {
		d.Add(time.Date(2022, time.January, day, 0, 0, 0, 0, time.UTC), "A", float64(day))
	}

	assert.Equal(t, 31, d.Size())

	d.FilterByRange(time.Time{}, time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, 31, d.Size())

	d.FilterByRange(time.Time{}, time.Date(2022, time.January, 30, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, 30, d.Size())

	d.FilterByRange(time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC), time.Time{})
	assert.Equal(t, 29, d.Size())

	d.FilterByRange(time.Date(2022, time.January, 8, 0, 0, 0, 0, time.UTC), time.Date(2022, time.January, 14, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, 7, d.Size())

	values, ok := d.GetValues("A")
	require.True(t, ok)
	assert.Equal(t, []float64{8, 9, 10, 11, 12, 13, 14}, values)

}

func BenchmarkDataset_FilterByRange(b *testing.B) {
	d := datasets.New()
	timestamp := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	for day := 0; day < 5*365; day++ {
		d.Add(timestamp, "A", float64(day))
		timestamp = timestamp.Add(24 * time.Hour)
	}

	b.ResetTimer()

	start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	stop := timestamp
	for i := 0; i < b.N; i++ {
		d.FilterByRange(start, stop)
		start = start.Add(12 * time.Hour)
		stop = stop.Add(-12 * time.Hour)
	}
}

func TestDataset_Accumulate(t *testing.T) {
	d := datasets.New()
	assert.NotNil(t, d)

	for day := 1; day < 32; day++ {
		d.Add(time.Date(2022, time.January, day, 0, 0, 0, 0, time.UTC), "A", 1.0)
	}

	d.Accumulate()

	values, ok := d.GetValues("A")
	require.True(t, ok)
	expected := 1.0
	for index, value := range values {
		require.Equal(t, expected, value, index)
		expected++
	}
}

func TestDataset_Copy(t *testing.T) {
	d := datasets.New()
	assert.NotNil(t, d)

	for day := 1; day < 5; day++ {
		ts := time.Date(2022, time.January, day, 0, 0, 0, 0, time.UTC)
		d.Add(ts, "A", 1.0)
	}

	clone := d.Copy()

	d.Accumulate()

	values, ok := d.GetValues("A")
	require.True(t, ok)
	assert.Equal(t, []float64{1, 2, 3, 4}, values)

	values, ok = clone.GetValues("A")
	require.True(t, ok)
	assert.Equal(t, []float64{1, 1, 1, 1}, values)
}
