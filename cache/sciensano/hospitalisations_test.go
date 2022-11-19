package sciensano_test

import (
	"encoding/json"
	"github.com/clambin/sciensano/cache/sciensano"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHospitalisations_Unmarshal(t *testing.T) {
	f, err := os.Open(filepath.Join("input", "hospitalisations.json"))
	require.NoError(t, err)

	var input sciensano.Hospitalisations
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}

func TestHospitalisations_Summarize(t *testing.T) {
	hospitalisations := makeHospitalisations(1)

	if *update {
		f, err := os.OpenFile(filepath.Join("input", "hospitalisations.json"), os.O_TRUNC|os.O_WRONLY, 0644)
		require.NoError(t, err)
		decoder := json.NewEncoder(f)
		decoder.SetIndent("", "  ")
		err = decoder.Encode(hospitalisations)
		require.NoError(t, err)
		_ = f.Close()
		t.Log("updated")
	}

	// TODO
}

func makeHospitalisations(count int) sciensano.Hospitalisations {
	return makeResponse[sciensano.Hospitalisation](count, func(timestamp time.Time, region, province, ageGroup, _ string, _ sciensano.DoseType) sciensano.Hospitalisation {
		return sciensano.Hospitalisation{
			TimeStamp:   sciensano.TimeStamp{Time: timestamp},
			Province:    province,
			Region:      region,
			TotalIn:     10,
			TotalInICU:  3,
			TotalInResp: 2,
			TotalInECMO: 1,
		}
	})
}
