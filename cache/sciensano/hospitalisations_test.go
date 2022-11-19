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

func TestHospitalisations_Unmarshal(t *testing.T) {
	f, err := os.Open(filepath.Join("input", "hospitalisations.json"))
	require.NoError(t, err)

	var input sciensano.Hospitalisations
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}
