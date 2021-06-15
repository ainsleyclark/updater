// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package archive

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/ainsleyclark/updater/pkg/fileio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Zip struct {
	Path string
	reader *zip.ReadCloser
}

func (z *Zip) Open() error {
	if !fileio.FileExists(z.Path) {
		return errors.New("no zip file exists with the path: "+z.Path)
	}
	reader, err := zip.OpenReader(z.Path)
	if err != nil {
		return err
	}
	z.reader = reader
	return nil
}

func (z *Zip) Close() error {
	if z.reader == nil {
		return nil
	}
	return z.reader.Close()
}

func (z *Zip) Unzip(dest string) ([]string, error) {
	var filenames []string

	for _, f := range z.reader.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		// Make folder if it's a directory
		if f.FileInfo().IsDir() {

			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm);
		if err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}