package table_test

import (
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/reporter/table"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

var input = []apiclient.APIResponse{
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

func TestNewDataframeFromAPIResponse(t *testing.T) {
	d := table.NewFromAPIResponse(makeBigResponse())

	timestamps := d.GetTimestamps()
	require.Len(t, timestamps, 720)
	assert.Equal(t, time.Date(2021, time.December, 20, 0, 0, 0, 0, time.UTC), timestamps[719])

	assert.Len(t, d.GetColumns(), 3)

	aValues, ok := d.GetFloatValues("Value1")
	require.True(t, ok)
	require.Len(t, aValues, 720)
	assert.Equal(t, 1.0, aValues[719])

	bValues, ok := d.GetFloatValues("Value2")
	require.True(t, ok)
	require.Len(t, bValues, 720)
	assert.Equal(t, 2.0, bValues[719])
}

func TestNewGroupedDataframeFromAPIResponse(t *testing.T) {
	d := table.NewGroupedFromAPIResponse(input, 1)

	timestamps, found := d.GetTimeValues("time")
	require.True(t, found)
	assert.Equal(t, []time.Time{
		time.Date(2022, time.February, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.February, 3, 0, 0, 0, 0, time.UTC),
	}, timestamps)

	assert.Len(t, d.Frame.Fields, 3)

	expectedValues := map[string][]float64{
		"A": {4, 5},
		"B": {2, 0},
	}

	for col, values := range expectedValues {
		v, found := d.GetFloatValues(col)
		require.True(t, found)
		assert.Equal(t, values, v)
	}
}

func makeBigResponse() (response []apiclient.APIResponse) {
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 720; i++ {
		response = append(response, &testMultiValueResponse{
			timestamp: timestamp,
			value1:    1,
			value2:    2,
		})
		timestamp = timestamp.Add(24 * time.Hour)
	}
	return
}

func BenchmarkNewFromAPIResponse(b *testing.B) {
	bigResponse := makeBigResponse()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = table.NewFromAPIResponse(bigResponse)
	}
}

func BenchmarkNewGroupedDataframeNewGroupedFromAPIResponse(b *testing.B) {
	var r []apiclient.APIResponse
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 720; i++ {
		for c := 0; c < 9; c++ {
			for d := 0; d < 12; d++ {
				r = append(r, &testResponse{
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
		_ = table.NewGroupedFromAPIResponse(r, 1)
	}
}

type testMultiValueResponse struct {
	timestamp time.Time
	value1    float64
	value2    float64
}

type testResponse struct {
	timestamp time.Time
	group     string
	value     float64
}

var _ apiclient.APIResponse = &testResponse{}

func (t testResponse) GetTimestamp() time.Time {
	return t.timestamp
}

func (t testResponse) GetGroupFieldValue(groupField apiclient.GroupField) (value string) {
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

var _ apiclient.APIResponse = &testMultiValueResponse{}

func (t testMultiValueResponse) GetTimestamp() time.Time {
	return t.timestamp
}

func (t testMultiValueResponse) GetGroupFieldValue(_ apiclient.GroupField) (value string) {
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
