package demographics

import (
	"fmt"
	"github.com/clambin/sciensano/metrics"
	"io"
	"net/http"
	"os"
	"path"
)

const demographicsURL = "https://statbel.fgov.be/sites/default/files/files/opendata/bevolking%20naar%20woonplaats%2C%20nationaliteit%20burgelijke%20staat%20%2C%20leeftijd%20en%20geslacht/TF_SOC_POP_STRUCT_2021.zip"

// DataFile represents a demographics data file
type DataFile struct {
	// TempDirectory specifies the directory to use for temporary files. Uses system-specified tempdir is left blank
	TempDirectory string
	// URL
	URL string

	tempdir  string
	filename string
}

// Download downloads a demographics data file to disk
func (datafile *DataFile) Download() (err error) {
	datafile.tempdir, err = datafile.makeTempDir()

	if err != nil {
		return
	}

	zipFile := path.Join(datafile.tempdir, "demographics.zip")

	err = datafile.get(zipFile)

	if err != nil {
		return
	}

	var files []string
	files, err = unzip(zipFile, datafile.tempdir)

	if err != nil {
		return
	}

	for _, file := range files {
		if path.Base(file) == "TF_SOC_POP_STRUCT_2021.txt" {
			datafile.filename = file
			return nil
		}
	}

	return fmt.Errorf("could not find population file in archive")

}

func (datafile *DataFile) makeTempDir() (name string, err error) {
	tempDir := datafile.TempDirectory
	if tempDir == "" {
		tempDir = os.TempDir()
	}
	return os.MkdirTemp(tempDir, "demographics")
}

func (datafile *DataFile) get(filename string) (err error) {
	timer := metrics.NewTimerMetric("demographics")

	url := datafile.URL
	if url == "" {
		url = demographicsURL
	}

	var resp *http.Response
	if resp, err = http.Get(url); err == nil {
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("server returned %s", resp.Status)
		} else {
			// Create the file
			var out *os.File
			if out, err = os.Create(filename); err == nil {
				// Write the body to file
				_, err = io.Copy(out, resp.Body)
				_ = out.Close()
			}
		}
		_ = resp.Body.Close()
	}

	timer.Report(err == nil)
	return
}

// Remove deletes the downloaded demographics data file from disk
func (datafile *DataFile) Remove() {
	_ = os.RemoveAll(datafile.tempdir)
}
