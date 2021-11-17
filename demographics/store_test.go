package demographics_test

import (
	"context"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/demographics/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	testServer := fake.New("")
	defer testServer.Close()

	store := demographics.Store{
		AgeBrackets: demographics.DefaultAgeBrackets,
		URL:         testServer.URL(),
	}

	store.Update()
	require.Len(t, store.Stats(), 2)

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

	stats := store.Stats()
	require.Len(t, stats, 2)
	require.Contains(t, stats, "Regions")
	assert.Equal(t, len(regions), stats["Regions"])
	require.Contains(t, stats, "AgeBrackets")
	assert.Equal(t, len(brackets), stats["AgeBrackets"])
}

func TestStore_Server_Failure(t *testing.T) {
	testServer := fake.New("")

	store := demographics.Store{
		AgeBrackets: demographics.DefaultAgeBrackets,
		URL:         testServer.URL(),
	}

	testServer.Close()
	_, ok := store.GetByRegion("Flanders")
	require.False(t, ok)
}

func TestStore_AutoRefresh(t *testing.T) {
	testServer := fake.New("")
	defer testServer.Close()

	store := demographics.Store{
		AgeBrackets: demographics.DefaultAgeBrackets,
		URL:         testServer.URL(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		store.AutoRefresh(ctx, 50*time.Millisecond)
		wg.Done()
	}()

	require.Eventually(t, func() bool {
		return len(store.Stats()) == 2
	}, 5*time.Second, 20*time.Millisecond)

	time.Sleep(200 * time.Millisecond)

	cancel()
	wg.Wait()
}

func BenchmarkStore(b *testing.B) {
	testServer := fake.New("../data/big_demographics.zip")
	defer testServer.Close()

	store := demographics.Store{
		AgeBrackets: demographics.DefaultAgeBrackets,
		URL:         testServer.URL(),
	}
	store.Update()
	_ = store.GetRegions()
}
