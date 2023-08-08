package testutil_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/clambin/sciensano/internal/sciensano/testutil"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "update sciensano reference responses")

func TestMain(m *testing.M) {
	flag.Parse()
	if *update {
		updateReferenceFiles()
	}
	m.Run()
}

func TestReferenceData(t *testing.T) {
	assert.NotEmpty(t, testutil.Cases())
	assert.NotEmpty(t, testutil.TestResults())
	assert.NotEmpty(t, testutil.Mortalities())
	assert.NotEmpty(t, testutil.Hospitalisations())
	assert.NotEmpty(t, testutil.Vaccinations())
}

func updateReferenceFiles() {
	updateReferenceFile[*sciensano.Case]("https://epistat.sciensano.be/Data/COVID19BE_CASES_AGESEX.json", "cases.json")
	updateReferenceFile[*sciensano.Hospitalisation]("https://epistat.sciensano.be/Data/COVID19BE_HOSP.json", "hospitalisations.json")
	updateReferenceFile[*sciensano.Mortality]("https://epistat.sciensano.be/Data/COVID19BE_MORT.json", "mortalities.json")
	updateReferenceFile[*sciensano.TestResult]("https://epistat.sciensano.be/Data/COVID19BE_tests.json", "testResults.json")
	updateReferenceFile[*sciensano.Vaccination]("https://epistat.sciensano.be/Data/COVID19BE_VACC.json", "vaccinations.json")
}

func updateReferenceFile[T any](source, filename string) {
	content, err := http.Get(source)
	if err != nil {
		panic(err)
	}
	defer content.Body.Close()

	filtered, err := filterReferenceFile[T](content.Body)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(path.Join("testdata", filename), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.Copy(f, filtered)
	if err != nil {
		panic(err)
	}
}

func filterReferenceFile[T any](body io.Reader) (io.ReadWriter, error) {
	var records []T
	if err := json.NewDecoder(body).Decode(&records); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	minTimestamp := time.Now().Add(-365 * 24 * time.Hour)

	var filtered []T

	for _, record := range records {
		if getRecordTimestamp(record).After(minTimestamp) {
			filtered = append(filtered, record)
		}
	}

	var out bytes.Buffer
	if err := json.NewEncoder(&out).Encode(filtered); err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}
	return &out, nil
}

func getRecordTimestamp(record any) time.Time {
	switch r := record.(type) {
	case *sciensano.Case:
		return r.TimeStamp.Time
	case *sciensano.Hospitalisation:
		return r.TimeStamp.Time
	case *sciensano.Mortality:
		return r.TimeStamp.Time
	case *sciensano.TestResult:
		return r.TimeStamp.Time
	case *sciensano.Vaccination:
		return r.TimeStamp.Time
	}
	panic(fmt.Errorf("invalid type: %v", record))
}

func BenchmarkCases(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = testutil.Cases()
	}
}
