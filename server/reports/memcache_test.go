package reports_test

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/clambin/go-common/tabulator"
	"github.com/clambin/sciensano/server/reports"
	"github.com/clambin/sciensano/server/reports/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMemCache(t *testing.T) {
	const encodedItem = `{"Timestamps":["2023-07-29T00:00:00Z","2023-07-30T00:00:00Z","2023-07-31T00:00:00Z"],"Columns":["A","B","C"],"Data":[[1,0,0],[0,2,0],[0,0,3]]}`

	m := mocks.NewMemCacheClient(t)
	m.On("Get", "foo").Return(nil, memcache.ErrCacheMiss).Once()
	m.On("Get", "foo").Return(&memcache.Item{Value: []byte(encodedItem)}, nil)
	m.On("Set", mock.AnythingOfType("*memcache.Item")).Return(nil).Once()

	c := reports.ReportCache{Cache: reports.NewMemCache(m, 15*time.Minute)}
	refTable, _ := createSimpleDataSet()

	var updates int
	for i := 0; i < 10; i++ {
		fromCache, err := c.MaybeGenerate("foo", func() (*tabulator.Tabulator, error) {
			updates++
			return createSimpleDataSet()
		})
		require.NoError(t, err)
		assert.Equal(t, refTable.Size(), fromCache.Size(), i)
		assert.Equal(t, refTable.GetColumns(), fromCache.GetColumns(), i)
	}

	assert.Equal(t, 1, updates)
}

func BenchmarkMemCache(b *testing.B) {
	c := reports.ReportCache{Cache: reports.NewMemCache(memcache.New("localhost:11211"), 15*time.Minute)}
	_, err := c.MaybeGenerate("foo", createBigDataSet)
	require.NoError(b, err)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = c.MaybeGenerate("foo", createBigDataSet)
		if err != nil {
			b.Fatal()
		}
	}
}
