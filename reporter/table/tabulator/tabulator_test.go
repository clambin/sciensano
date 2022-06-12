package tabulator_test

import (
	"github.com/clambin/sciensano/reporter/table/tabulator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

//var update = flag.Bool("update", false, "update .golden files")

func TestTabulator(t *testing.T) {
	d := tabulator.New("B")
	assert.NotNil(t, d)

	added := d.Add(time.Now(), "foo", 1.0)
	assert.False(t, added)

	d.RegisterColumn("A", "B")

	for day := 1; day < 5; day++ {
		added = d.Add(time.Date(2022, time.January, day, 0, 0, 0, 0, time.UTC), "A", float64(day))
		assert.True(t, added)
		added = d.Add(time.Date(2022, time.January, day, 0, 0, 0, 0, time.UTC), "B", -float64(day))
		assert.True(t, added)
	}

	assert.Equal(t, 4, d.Size())
	assert.Equal(t, []string{"A", "B"}, d.GetColumns())
	assert.Equal(t, []time.Time{
		time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.January, 4, 0, 0, 0, 0, time.UTC),
	}, d.GetTimestamps())

	values, ok := d.GetValues("A")
	require.True(t, ok)
	assert.Equal(t, []float64{1, 2, 3, 4}, values)

	values, ok = d.GetValues("B")
	require.True(t, ok)
	assert.Equal(t, []float64{-1, -2, -3, -4}, values)

	d.RegisterColumn("C")
	values, ok = d.GetValues("C")
	require.True(t, ok)
	assert.Equal(t, []float64{0, 0, 0, 0}, values)
}

func BenchmarkTabulator_Add(b *testing.B) {
	d := tabulator.New()
	d.RegisterColumn("A")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timestamp := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
		for day := 0; day < 5*365; day++ {
			d.Add(timestamp, "A", float64(day))
			timestamp = timestamp.Add(24 * time.Hour)
		}
	}
}

func BenchmarkTabulator_GetColumns(b *testing.B) {
	d := tabulator.New()
	d.RegisterColumn("A")
	timestamp := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	for day := 0; day < 5*365; day++ {
		d.Add(timestamp, "A", float64(day))
		timestamp = timestamp.Add(24 * time.Hour)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.GetColumns()
	}
}

func BenchmarkTabulator_GetTimestamps(b *testing.B) {
	d := tabulator.New()
	timestamp := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	for day := 0; day < 5*365; day++ {
		d.Add(timestamp, "A", float64(day))
		timestamp = timestamp.Add(24 * time.Hour)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.GetTimestamps()
	}
}

func BenchmarkTabulator_GetValues(b *testing.B) {
	d := tabulator.New()
	timestamp := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	for day := 0; day < 5*365; day++ {
		d.Add(timestamp, "A", float64(day))
		timestamp = timestamp.Add(24 * time.Hour)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.GetValues("A")
	}
}
