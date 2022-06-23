package sciensano_test

import (
	"encoding/json"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/sciensano"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAPIVaccinationsResponses(t *testing.T) {
	gp := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		body, err := easyjson.Marshal(vaccinationResponses)
		require.NoError(t, err)

		err = os.WriteFile(gp, body, 0644)
		require.NoError(t, err)
	}

	body, err := os.ReadFile(gp)
	require.NoError(t, err)

	var output sciensano.APIVaccinationsResponses
	err = easyjson.Unmarshal(body, &output)
	require.NoError(t, err)
	require.Len(t, output, len(casesResponses))
}

func TestAPIVaccinationsResponse_Attributes(t *testing.T) {
	timestamps := []time.Time{
		time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.June, 18, 0, 0, 0, 0, time.UTC),
	}
	groups := []string{"55-54", "45-54"}
	doses := [][]float64{{1, 0, 0, 0}, {0, 2, 0, 0}}
	for idx, entry := range vaccinationResponses {
		assert.Equal(t, timestamps[idx], entry.GetTimestamp(), idx)
		assert.Equal(t, []string{"partial", "full", "singledose", "booster"}, entry.GetAttributeNames())
		assert.Equal(t, doses[idx], entry.GetAttributeValues())
		assert.Equal(t, float64(idx+1), entry.GetTotalValue())
		assert.Equal(t, "Flanders", entry.GetGroupFieldValue(apiclient.GroupByRegion))
		assert.Equal(t, groups[idx], entry.GetGroupFieldValue(apiclient.GroupByAgeGroup))
		assert.Equal(t, "A", entry.GetGroupFieldValue(apiclient.GroupByManufacturer))
	}
}

var (
	vaccinationResponses = sciensano.APIVaccinationsResponses{
		&sciensano.APIVaccinationsResponse{
			TimeStamp:    sciensano.TimeStamp{Time: time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "A",
			Region:       "Flanders",
			AgeGroup:     "55-54",
			Dose:         "A",
			Count:        1,
		},
		&sciensano.APIVaccinationsResponse{
			TimeStamp:    sciensano.TimeStamp{Time: time.Date(2022, time.June, 18, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "A",
			Region:       "Flanders",
			AgeGroup:     "45-54",
			Dose:         "B",
			Count:        2,
		},
	}
)

func BenchmarkAPIVaccinationsResponses_UnmarshalEasyJSON(b *testing.B) {
	var body []byte
	var err error
	if body, err = os.ReadFile("../../data/vaccinations.json"); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var response sciensano.APIVaccinationsResponses
		if err = easyjson.Unmarshal(body, &response); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAPIVaccinationsResponses_UnmarshalJSON(b *testing.B) {
	type vaccinationEntry struct {
		TimeStamp    sciensano.TimeStamp `json:"DATE"`
		Manufacturer string              `json:"BRAND"`
		Region       string              `json:"REGION"`
		AgeGroup     string              `json:"AGEGROUP"`
		Gender       string              `json:"SEX"`
		Dose         string              `json:"DOSE"`
		Count        int                 `json:"COUNT"`
	}

	var body []byte
	var err error
	if body, err = os.ReadFile("../../data/vaccinations.json"); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var response []vaccinationEntry
		if err = json.Unmarshal(body, &response); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAPIVaccinationsResponses_UnmarshalJSONIter(b *testing.B) {
	type vaccinationEntry struct {
		TimeStamp    sciensano.TimeStamp `json:"DATE"`
		Manufacturer string              `json:"BRAND"`
		Region       string              `json:"REGION"`
		AgeGroup     string              `json:"AGEGROUP"`
		Gender       string              `json:"SEX"`
		Dose         string              `json:"DOSE"`
		Count        int                 `json:"COUNT"`
	}

	var body []byte
	var err error
	if body, err = os.ReadFile("../../data/vaccinations.json"); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var response []vaccinationEntry
		if err = jsoniter.Unmarshal(body, &response); err != nil {
			b.Fatal(err)
		}
	}
}
