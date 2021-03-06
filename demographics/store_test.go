package demographics_test

import (
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/demographics/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	testServer := fake.New("")
	defer testServer.Close()

	store := demographics.Store{
		Retention:   time.Hour,
		AgeBrackets: demographics.DefaultAgeBrackets,
		URL:         testServer.URL(),
	}

	brackets := store.GetAgeBrackets()
	require.Len(t, brackets, len(demographics.DefaultAgeBrackets)+1)

	for _, bracket := range brackets {
		population, ok := store.GetByAge(bracket)
		require.True(t, ok, bracket)
		assert.NotZero(t, population)
	}

	figures := store.GetAgeGroupFigures()
	assert.Len(t, figures, len(demographics.DefaultAgeBrackets)+1)

	for _, bracket := range brackets {
		assert.Contains(t, figures, bracket.String())
	}

	regions := store.GetRegions()
	require.Len(t, regions, 3)

	for _, region := range store.GetRegions() {
		population, ok := store.GetByRegion(region)
		require.True(t, ok)
		assert.NotZero(t, population)
	}

	figures = store.GetRegionFigures()
	for _, region := range regions {
		assert.Contains(t, figures, region)
	}
}

func TestStore_Server_Failure(t *testing.T) {
	testServer := fake.New("")

	store := demographics.Store{
		Retention:   time.Hour,
		AgeBrackets: demographics.DefaultAgeBrackets,
		URL:         testServer.URL(),
	}

	testServer.Close()
	_, ok := store.GetByRegion("Flanders")
	require.False(t, ok)
}

func BenchmarkStore(b *testing.B) {
	testServer := fake.New("../data/big_demographics.zip")

	store := demographics.Store{
		Retention:   time.Hour,
		AgeBrackets: demographics.DefaultAgeBrackets,
		URL:         testServer.URL(),
	}

	_ = store.GetRegions()
}
