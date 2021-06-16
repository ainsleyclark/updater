// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zip

import (
	"github.com/ainsleyclark/updater/tests"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestZip_Copy(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	zipDir := filepath.Join(wd, "../../..", tests.DataPath, "zips")

	tt := map[string]struct {
		src  string
		dest string
		want interface{}
	}{
		"Success": {
			filepath.Join(zipDir, "image.zip"),
			t.TempDir(),
			nil,
		},
		"Wrong File Path": {
			"wrong",
			"",
			"open wrong: no such file or directory",
		},
		"Illegal file path": {
			filepath.Join(zipDir, "image.zip"),
			".",
			"illegal file path",
		},
		"With Folder": {
			filepath.Join(zipDir, "with-folder.zip"),
			t.TempDir(),
			nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			zip := Zip{Path: test.src}
			got := zip.Copy(test.dest)
			if got != nil {
				require.Contains(t, got.Error(), test.want)
				return
			}
			require.Equal(t, test.want, got)
		})
	}
}
