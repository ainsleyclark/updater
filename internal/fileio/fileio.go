// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fileio

import (
	"github.com/kardianos/osext"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	// TempDirName is the name used to store any temporary
	// directories used by the updater.
	TempDirName = "verbis-updater"
)

// Exists checks to see if a file or directory exists
// by the given path.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDirectory returns true if the source path is a
// directory. Returns false if it is a file.
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
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

// SplitPaths
func SplitPaths(base, target string) (string, error) {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return "", err
	}
	if rel == "." {
		return "", nil
	}
	// TODO: Check on Windows
	return strings.ReplaceAll(rel, ".."+string(os.PathSeparator), ""), nil
}
