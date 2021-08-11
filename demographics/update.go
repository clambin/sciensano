package demographics

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

func (server *Server) update(ctx context.Context) (err error) {
	server.lock.Lock()
	defer server.lock.Unlock()

	var tempDir, filename string

	tempDir, err = server.makeTempDir()

	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %s", err.Error())
	}

	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	start := time.Now()
	filename, err = server.download(ctx, tempDir)

	if err != nil {
		return fmt.Errorf("unable to download demographics data: %s", err.Error())
	}

	// process file here
	log.Debug(filename)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		var localErr error
		server.byAge, localErr = server.parseByAge(filename)

		if localErr != nil {
			log.WithError(localErr).Fatal("unable to parse demographics file")
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		var localErr error
		server.byRegion, localErr = server.parseByRegion(filename)

		if localErr != nil {
			log.WithError(localErr).Fatal("unable to parse demographics file")
		}
		wg.Done()
	}()

	wg.Wait()
	log.Infof("loaded demographics in %s", time.Now().Sub(start))
	return
}

const demographicsURL = "https://statbel.fgov.be/sites/default/files/files/opendata/bevolking%20naar%20woonplaats%2C%20nationaliteit%20burgelijke%20staat%20%2C%20leeftijd%20en%20geslacht/TF_SOC_POP_STRUCT_2021.zip"

func (server *Server) download(ctx context.Context, tempDir string) (filename string, err error) {
	zipFile := path.Join(tempDir, "demographics.zip")
	err = server.get(ctx, zipFile)

	if err != nil {
		return
	}

	var files []string
	files, err = unzip(zipFile, tempDir)

	if err != nil {
		return
	}

	for _, file := range files {
		if path.Base(file) == "TF_SOC_POP_STRUCT_2021.txt" {
			return file, nil
		}
	}

	return "", fmt.Errorf("could not find population file in archive")
}

func (server *Server) makeTempDir() (path string, err error) {
	tempDir := server.TempDirectory
	if tempDir == "" {
		tempDir = os.TempDir()
	}
	return os.MkdirTemp(tempDir, "demographics")
}

func (server *Server) get(ctx context.Context, filename string) (err error) {
	url := server.URL
	if url == "" {
		url = demographicsURL
	}

	var req *http.Request
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	var resp *http.Response
	resp, err = server.HTTPClient.Do(req)

	if err != nil {
		return
	}

	defer func(body io.ReadCloser) {
		_ = resp.Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %s", resp.Status)
	}

	// Create the file
	var out *os.File
	out, err = os.Create(filename)

	if err == nil {
		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		_ = out.Close()
	}

	return
}
