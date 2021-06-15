// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fileio

import (
	"github.com/kardianos/osext"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// TempDirName is the name used to store any temporary
	// directories used by the updater.
	TempDirName = "verbis-updater"
)

// FileExists checks to see if a file or directory exists
// by the given path.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Executable finds the current Executable and Folder,
// returning an error if it could not be found.
func Executable() (string, error) {
	return osext.Executable()
}

// TempDirectory creates a temporary directory and returns
// a path upon success.
func TempDirectory() (string, error) {
	return ioutil.TempDir("", TempDirName)
}

type Paths struct {
	Base           string
	ExecutableName string
}

func GetPaths() (*Paths, error) {
	exec, err := Executable()
	if err != nil {
		return nil, err
	}
	return &Paths{
		Base:           filepath.Base(exec),
		ExecutableName: exec,
	}, nil
}
