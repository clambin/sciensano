package demographics_test

import (
	"context"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/demographics/mock"
	"github.com/stretchr/testify/assert"
	"math"
	"sync"
	"testing"
	"time"
)

func TestServer_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	testServer := mock.New("")
	defer testServer.Close()

	server := demographics.New()
	server.URL = testServer.URL()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := server.Run(ctx, 24*time.Hour)
		assert.NoError(t, err)
		wg.Done()
	}()

	assert.Eventually(t, server.AvailableData, 1000*time.Millisecond, 10*time.Millisecond)

	count, ok := server.GetByAge(demographics.Bracket{Low: 0, High: 11})
	assert.True(t, ok)
	assert.NotZero(t, count)

	brackets := server.GetAgeBrackets()
	assert.Contains(t, brackets, demographics.Bracket{Low: 0, High: 11})
	assert.Contains(t, brackets, demographics.Bracket{Low: 12, High: 15})
	assert.Contains(t, brackets, demographics.Bracket{Low: 16, High: 17})
	assert.Contains(t, brackets, demographics.Bracket{Low: 18, High: 24})
	assert.Contains(t, brackets, demographics.Bracket{Low: 25, High: 34})
	assert.Contains(t, brackets, demographics.Bracket{Low: 35, High: 44})
	assert.Contains(t, brackets, demographics.Bracket{Low: 45, High: 54})
	assert.Contains(t, brackets, demographics.Bracket{Low: 55, High: 64})
	assert.Contains(t, brackets, demographics.Bracket{Low: 65, High: 74})
	assert.Contains(t, brackets, demographics.Bracket{Low: 75, High: 84})
	assert.Contains(t, brackets, demographics.Bracket{Low: 85, High: math.Inf(+1)})

	regions := server.GetRegions()
	assert.Contains(t, regions, "Flanders")
	assert.Contains(t, regions, "Wallonia")
	assert.Contains(t, regions, "Brussels")

	cancel()
	wg.Wait()
}

func BenchmarkServer_Run(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testServer := mock.New("../data/big_demographics.zip")
	defer testServer.Close()

	server := demographics.New()
	server.URL = testServer.URL()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := server.Run(ctx, 24*time.Hour)
		assert.NoError(b, err)
		wg.Done()
	}()

	assert.Eventually(b, server.AvailableData, 1000*time.Millisecond, 10*time.Millisecond)
}
