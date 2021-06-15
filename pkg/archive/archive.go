// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package archive

import (
	"fmt"
	"github.com/ainsleyclark/updater/internal/fileio"
	"github.com/ainsleyclark/updater/pkg/archive/zip"
	"strings"
)

type Archive interface {
	Copy(dest string) error
}

func New(src string) (Archive, error) {
	if !fileio.FileExists(src) {
		return nil, fmt.Errorf("file not found with the path: %s", src)
	}
	if strings.HasSuffix(src, ".zip") {
		return &zip.Zip{Path: src}, nil
	}
	return nil, fmt.Errorf("no archiver found")
}