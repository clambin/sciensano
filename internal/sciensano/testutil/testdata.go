package testutil

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/clambin/sciensano/internal/sciensano"
	"path"
)

//go:embed testdata/*
var testFiles embed.FS

func Cases() sciensano.Cases {
	return getTestData[sciensano.Cases]("cases.json")
}

func TestResults() sciensano.TestResults {
	return getTestData[sciensano.TestResults]("testResults.json")
}

func Mortalities() sciensano.Mortalities {
	return getTestData[sciensano.Mortalities]("mortalities.json")
}

func Hospitalisations() sciensano.Hospitalisations {
	return getTestData[sciensano.Hospitalisations]("hospitalisations.json")
}

func Vaccinations() sciensano.Vaccinations {
	return getTestData[sciensano.Vaccinations]("vaccinations.json")
}

func getTestData[T any](filename string) T {
	f, err := testFiles.Open(path.Join("testdata", filename))
	if err != nil {
		panic(fmt.Errorf("testdata %s: %w", filename, err))
	}
	var records T
	if err = json.NewDecoder(f).Decode(&records); err != nil {
		panic(fmt.Errorf("testdata %s unmarshal: %w", filename, err))
	}
	return records
}
