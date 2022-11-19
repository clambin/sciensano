package sciensano_test

import (
	"flag"
	"fmt"
	"github.com/clambin/sciensano/cache/sciensano"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

var (
	regions       = []string{"", "Brussels", "Flanders", "Ostbelgien", "Wallonia"}
	ageGroups     = []string{"00-04", "05-11", "12-15", "16-17", "18-24", "25-34", "35-44", "45-54", "55-64", "65-74", "75-84", "85+"}
	manufacturers = []string{"AstraZeneca-Oxford", "Johnson&Johnson", "Moderna", "Novavax", "Other", "Pfizer-BioNTech"}
	doses         = []sciensano.DoseType{sciensano.Partial, sciensano.Full, sciensano.SingleDose, sciensano.Booster, sciensano.Booster2, sciensano.Booster3}
	update        = flag.Bool("update", false, "update input files")
)

func TestMain(m *testing.M) {
	// FIXME: update is always the default value?
	if *update {
		if err := updateInputFiles(); err != nil {
			panic(err)
		}
	}
	m.Run()
}

func updateInputFiles() error {
	for name, route := range sciensano.Routes {
		resp, err := http.Get(sciensano.BaseURL + route)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("%s: %s", name, resp.Status)
		}

		filename := filepath.Join("input", name+".json")
		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("open %s: %w", filename, err)
		}

		if _, err = io.Copy(f, resp.Body); err != nil {
			return fmt.Errorf("write %s: %w", filename, err)
		}
		_ = f.Close()
	}
	return nil
}
