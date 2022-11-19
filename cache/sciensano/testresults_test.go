package sciensano_test

import (
	"encoding/json"
	"github.com/clambin/sciensano/cache/sciensano"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestTestResults_Unmarshal(t *testing.T) {
	f, err := os.Open(filepath.Join("input", "testResults.json"))
	require.NoError(t, err)

	var input sciensano.TestResults
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}
