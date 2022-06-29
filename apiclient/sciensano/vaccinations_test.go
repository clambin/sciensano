package sciensano_test

import (
	"encoding/json"
	"fmt"
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
	responses := buildVaccinationResponses()

	gp := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		body, err := easyjson.Marshal(responses)
		require.NoError(t, err)

		err = os.WriteFile(gp, body, 0644)
		require.NoError(t, err)
	}

	body, err := os.ReadFile(gp)
	require.NoError(t, err)

	var output sciensano.APIVaccinationsResponses
	err = easyjson.Unmarshal(body, &output)
	require.NoError(t, err)
	require.Len(t, output, len(responses))
}

func TestAPIVaccinationsResponse_Attributes(t *testing.T) {
	for _, entry := range buildVaccinationResponses() {
		assert.Equal(t, time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC), entry.GetTimestamp())
		assert.Equal(t, []string{"partial", "full", "singledose", "booster", "booster2"}, entry.GetAttributeNames())

		var attributeValues []float64
		var totalValue float64
		switch entry.Dose {
		case sciensano.TypeVaccinationPartial:
			attributeValues = []float64{0, 0, 0, 0, 0}
			totalValue = 0.0
		case sciensano.TypeVaccinationFull:
			attributeValues = []float64{0, 1, 0, 0, 0}
			totalValue = 1.0
		case sciensano.TypeVaccinationSingle:
			attributeValues = []float64{0, 0, 2, 0, 0}
			totalValue = 2.0
		case sciensano.TypeVaccinationBooster:
			attributeValues = []float64{0, 0, 0, 3, 0}
			totalValue = 3.0
		case sciensano.TypeVaccinationBooster2:
			attributeValues = []float64{0, 0, 0, 0, 4}
			totalValue = 4.0
		default:
			t.Fatalf("unexpected dose value: %d", int(entry.Dose))
		}

		assert.Equal(t, attributeValues, entry.GetAttributeValues())
		assert.Equal(t, totalValue, entry.GetTotalValue())
		assert.Equal(t, "Flanders", entry.GetGroupFieldValue(apiclient.GroupByRegion))
		assert.Equal(t, "45-54", entry.GetGroupFieldValue(apiclient.GroupByAgeGroup))
		assert.Equal(t, "A", entry.GetGroupFieldValue(apiclient.GroupByManufacturer))
	}
}

func TestDoseType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		d        sciensano.DoseType
		body     string
		expected sciensano.DoseType
		wantErr  assert.ErrorAssertionFunc
	}{
		{body: `"A"`, expected: sciensano.TypeVaccinationPartial, wantErr: assert.NoError},
		{body: `"B"`, expected: sciensano.TypeVaccinationFull, wantErr: assert.NoError},
		{body: `"C"`, expected: sciensano.TypeVaccinationSingle, wantErr: assert.NoError},
		{body: `"E"`, expected: sciensano.TypeVaccinationBooster, wantErr: assert.NoError},
		{body: `"E2"`, expected: sciensano.TypeVaccinationBooster2, wantErr: assert.NoError},
		{body: `""`, wantErr: assert.Error},
		{body: `A`, wantErr: assert.Error},
	}
	for _, tt := range tests {
		tt.wantErr(t, tt.d.UnmarshalJSON([]byte(tt.body)), fmt.Sprintf("UnmarshalJSON(%v)", tt.body))
	}
}

func TestDoseType_MarshalJSON(t *testing.T) {
	tests := []struct {
		d        sciensano.DoseType
		wantBody string
		wantErr  assert.ErrorAssertionFunc
	}{
		{d: sciensano.DoseType(sciensano.TypeVaccinationPartial), wantBody: `"A"`, wantErr: assert.NoError},
		{d: sciensano.DoseType(sciensano.TypeVaccinationFull), wantBody: `"B"`, wantErr: assert.NoError},
		{d: sciensano.DoseType(sciensano.TypeVaccinationSingle), wantBody: `"C"`, wantErr: assert.NoError},
		{d: sciensano.DoseType(sciensano.TypeVaccinationBooster), wantBody: `"E"`, wantErr: assert.NoError},
		{d: sciensano.DoseType(sciensano.TypeVaccinationBooster2), wantBody: `"E2"`, wantErr: assert.NoError},
		{d: sciensano.DoseType(-1), wantErr: assert.Error},
	}

	for _, tt := range tests {
		gotBody, err := tt.d.MarshalJSON()
		if !tt.wantErr(t, err, "MarshalJSON()") {
			return
		}
		assert.Equalf(t, tt.wantBody, string(gotBody), "MarshalJSON()")
	}
}

func buildVaccinationResponses() (responses sciensano.APIVaccinationsResponses) {
	for _, doseType := range []sciensano.DoseType{
		sciensano.TypeVaccinationPartial,
		sciensano.TypeVaccinationFull,
		sciensano.TypeVaccinationSingle,
		sciensano.TypeVaccinationBooster,
		sciensano.TypeVaccinationBooster2,
	} {
		responses = append(responses, &sciensano.APIVaccinationsResponse{
			TimeStamp:    sciensano.TimeStamp{Time: time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "A",
			Region:       "Flanders",
			AgeGroup:     "45-54",
			Dose:         doseType,
			Count:        int(doseType),
		})
	}
	return
}

var (
	vaccinationResponses = sciensano.APIVaccinationsResponses{
		&sciensano.APIVaccinationsResponse{
			TimeStamp:    sciensano.TimeStamp{Time: time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "A",
			Region:       "Flanders",
			AgeGroup:     "55-54",
			Dose:         sciensano.TypeVaccinationPartial,
			Count:        1,
		},
		&sciensano.APIVaccinationsResponse{
			TimeStamp:    sciensano.TimeStamp{Time: time.Date(2022, time.June, 18, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "A",
			Region:       "Flanders",
			AgeGroup:     "45-54",
			Dose:         sciensano.TypeVaccinationFull,
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
		Dose         sciensano.DoseType  `json:"DOSE"`
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
		Dose         sciensano.DoseType  `json:"DOSE"`
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
		if err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, &response); err != nil {
			b.Fatal(err)
		}
	}
}
