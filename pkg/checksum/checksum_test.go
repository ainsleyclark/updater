// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checksum

import (
	"fmt"
	"github.com/ainsleyclark/updater/tests"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

const (
	ValidChecksumURL    = "https://github.com/cli/cli/releases/download/v1.11.0/gh_1.11.0_checksums.txt"
	ValidChecksumPath   = "gh_1.11.0_linux_amd64.tar.gz"
	InvalidChecksumPath = "gh_1.11.0_windows_amd64.zip"
)

func TestCompare(t *testing.T) {
	fi, err := tests.GetFileInfo()
	assert.NoError(t, err)

	tt := map[string]struct {
		url  string
		path string
		want error
	}{
		"Success": {
			ValidChecksumURL,
			filepath.Join(fi.TestDataPath, ValidChecksumPath),
			nil,
		},
		"Wrong URL": {
			"https://wrong.com",
			"",
			fmt.Errorf("EOF"),
		},
		"Wrong Path": {
			ValidChecksumURL,
			"wrong",
			ErrNoCheckSum,
		},
		"Mismatch": {
			ValidChecksumURL,
			filepath.Join(fi.TestDataPath, InvalidChecksumPath),
			ErrMismatch,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := Compare(test.url, test.path)
			if got != nil {
				assert.Contains(t, got.Error(), test.want.Error())
				return
			}
			assert.Equal(t, test.want, got)
		})
	}
}
