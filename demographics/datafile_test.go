package demographics_test

import (
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/demographics/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestDataFile(t *testing.T) {
	testServer := fake.New("../data/demographics.zip")
	defer testServer.Close()

	datafile := demographics.DataFile{
		URL: testServer.URL(),
	}

	err := datafile.Download()
	require.NoError(t, err)

	var byRegion map[string]int
	byRegion, err = datafile.ParseByRegion()
	require.NoError(t, err)
	assert.NotEmpty(t, byRegion)
	assert.NotEmpty(t, byRegion["Flanders"])
	assert.NotEmpty(t, byRegion["Wallonia"])
	assert.NotEmpty(t, byRegion["Brussels"])
	assert.Empty(t, byRegion["Ostbelgien"])

	var byAge map[demographics.Bracket]int
	byAge, err = datafile.ParseByAge([]float64{11, 25, 35})
	require.NoError(t, err)
	assert.NotEmpty(t, byAge)
	assert.NotEmpty(t, byAge[demographics.Bracket{Low: 0, High: 10}])
	assert.NotEmpty(t, byAge[demographics.Bracket{Low: 11, High: 24}])
	assert.NotEmpty(t, byAge[demographics.Bracket{Low: 25, High: 34}])
	assert.NotEmpty(t, byAge[demographics.Bracket{Low: 35, High: math.Inf(+1)}])

	datafile.Remove()
}
