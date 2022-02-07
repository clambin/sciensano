package demographics

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_Parser(t *testing.T) {
	byRegion, byAge, err := groupPopulation("../data/demographics.txt")
	require.NoError(t, err)
	require.Len(t, byRegion, 3)
	assert.Contains(t, byRegion, "Wallonia")
	assert.Contains(t, byRegion, "Flanders")
	assert.Contains(t, byRegion, "Brussels")
	assert.NotEmpty(t, byAge)
	assert.Contains(t, byAge, 52)
}

func BenchmarkStore_Parser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := groupPopulation("../data/TF_SOC_POP_STRUCT_2021.txt")
		if err != nil {
			b.Fatal(err)
		}
	}
}
