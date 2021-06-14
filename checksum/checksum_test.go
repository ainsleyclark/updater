// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checksum

import (
	"fmt"
	"github.com/ainsleyclark/updater/github"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)


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
