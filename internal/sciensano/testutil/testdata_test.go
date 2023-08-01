package testutil_test

import (
	"flag"
	"github.com/clambin/sciensano/internal/sciensano/testutil"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
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
	updateReferenceFile("https://epistat.sciensano.be/Data/COVID19BE_CASES_AGESEX.json", "cases.json")
	updateReferenceFile("https://epistat.sciensano.be/Data/COVID19BE_HOSP.json", "hospitalisations.json")
	updateReferenceFile("https://epistat.sciensano.be/Data/COVID19BE_tests.json", "testResults.json")
	updateReferenceFile("https://epistat.sciensano.be/Data/COVID19BE_MORT.json", "mortalities.json")
	updateReferenceFile("https://epistat.sciensano.be/Data/COVID19BE_VACC.json", "vaccinations.json")
}

func updateReferenceFile(source, filename string) {
	content, err := http.Get(source)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(path.Join("testdata", filename), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(f, content.Body)
	if err != nil {
		panic(err)
	}
}

func BenchmarkCases(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = testutil.Cases()
	}
}
