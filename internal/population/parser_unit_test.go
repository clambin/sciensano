package population

import (
	"archive/zip"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path"
	"testing"
)

var tmpDir string

func TestMain(m *testing.M) {
	if err := unzipTestFiles(); err != nil {
		panic(err)
	}

	m.Run()

	_ = os.RemoveAll(tmpDir)
	os.Exit(0)
}

func unzipTestFiles() (err error) {
	tmpDir, err = os.MkdirTemp("", "TestUnzip")
	if err != nil {
		panic(err)
	}

	for _, zipFile := range []string{"demographics.zip", "small_demographics.zip"} {
		var archive *zip.ReadCloser
		archive, err = zip.OpenReader(path.Join("testdata", zipFile))
		if err != nil {
			return fmt.Errorf("unable to open '%s': %w", zipFile, err)
		}
		for _, f := range archive.File {
			dstName := path.Join(tmpDir, f.Name)
			var dstFile *os.File
			dstFile, err = os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fmt.Errorf("could not create output file '%s': %w", dstName, err)
			}

			var fileInArchive io.ReadCloser
			fileInArchive, err = f.Open()
			if err != nil {
				return fmt.Errorf("could not open '%s' in archive '%s': $w", f.Name, zipFile)
			}

			if _, err = io.Copy(dstFile, fileInArchive); err != nil {
				return fmt.Errorf("failed to write '%s': %w", dstName, err)
			}

			_ = dstFile.Close()
			_ = fileInArchive.Close()
		}
	}
	return
}

func TestStore_Parser(t *testing.T) {
	byRegion, byAge, err := groupPopulation(path.Join(tmpDir, "demographics.txt"))
	require.NoError(t, err)
	require.Len(t, byRegion, 3)
	assert.Contains(t, byRegion, "Wallonia")
	assert.Contains(t, byRegion, "Flanders")
	assert.Contains(t, byRegion, "Brussels")
	assert.NotEmpty(t, byAge)
	assert.Contains(t, byAge, 52)
}

func BenchmarkStore_Parser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, err := groupPopulation(path.Join(tmpDir, "TF_SOC_POP_STRUCT_2021.txt"))
		if err != nil {
			b.Fatal(err)
		}
	}
}
