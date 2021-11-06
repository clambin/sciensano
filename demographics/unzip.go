package demographics

import (
	"archive/zip"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func unzip(infile, tmpdir string) (files []string, err error) {
	var reader *zip.ReadCloser
	reader, err = zip.OpenReader(infile)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = reader.Close()
	}()

	for _, f := range reader.File {
		filePath := path.Join(tmpdir, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(filePath, filepath.Clean(tmpdir)+string(os.PathSeparator)) {
			log.WithField("name", filePath).Warning("illegal f path. skipping")
			continue
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			_ = os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		files = append(files, filePath)

		// Make File
		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return
		}

		var outFile *os.File
		outFile, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return
		}

		var rc io.ReadCloser
		rc, err = f.Open()
		if err != nil {
			return
		}

		_, err = io.Copy(outFile, rc)

		// Close the f without defer to close before next iteration of loop
		_ = outFile.Close()
		_ = rc.Close()
	}

	return
}
