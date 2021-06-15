// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"fmt"
	"github.com/ainsleyclark/updater/pkg/github"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// UpdaterTestSuite defines the helper used for mail
// testing.
type UpdaterTestSuite struct {
	suite.Suite
	base     string
	testPath string
}

// Assert testing has begun.
func TestUpdater(t *testing.T) {
	suite.Run(t, new(UpdaterTestSuite))
}

// Assigns test base.
func (t *UpdaterTestSuite) SetupSuite() {
	wd, err := os.Getwd()
	t.NoError(err)
	t.base = filepath.Dir(wd)
	t.testPath = t.base + string(os.PathSeparator) + "testdata"
}

func TestUpdater_Update(t *testing.T) {
	u := Updater{
		Github: github.Repo{
			RepositoryURL: "https://github.com/ainsleyclark/verbis",
			ArchiveName:   fmt.Sprintf("verbis_%s_%s_%s.zip", "0.0.1", runtime.GOOS, runtime.GOARCH),
			ChecksumName:  "checksums.txt",
		},
		Files:                nil,
		Version:              "",
		BackupExtension:      "",
		RemoteExecutablePath: "verbis",
	}

	err := u.Update()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(err)
}

func (t *UpdaterTestSuite) TestFiles_Validate() {
	tt := map[string]struct {
		input string
		want  interface{}
	}{
		"Exists - File": {
			t.testPath + string(os.PathSeparator) + "gopher.png",
			nil,
		},
		"Exists - Folder": {
			t.testPath + string(os.PathSeparator) + "folder",
			nil,
		},
		"Not Exist": {
			"wrong",
			"no file or directory exists with the path",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			f := Files{File{LocalPath: test.input}}
			got := f.Validate()
			if got != nil {
				t.Contains(got.Error(), test.want)
				return
			}
			t.Equal(test.want, got)
		})
	}
}
