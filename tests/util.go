// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"os"
	"path/filepath"
)

const (
	ReleaseRepo = "github.com/ainsleyclark/verbis"
)

type Paths struct {
	BasePath string
	TestDataPath string
}

// GetFileInfo retrieves the paths relevant for testing.
func GetFileInfo() (*Paths, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	base := filepath.Dir(wd)

	return &Paths{
		BasePath:     base,
		TestDataPath: base + string(os.PathSeparator) + "testdata",
	}, nil
}