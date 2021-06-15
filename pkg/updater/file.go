// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"errors"
	fileio2 "github.com/ainsleyclark/updater/pkg/fileio"
)

type File struct {
	LocalPath  string
	RemotePath string
}

type Files []File

func (f Files) Validate() error {
	for _, file := range f {
		if !fileio2.FileExists(file.LocalPath) {
			return errors.New("no file or directory exists with the path: " + file.LocalPath)
		}
	}
	return nil
}