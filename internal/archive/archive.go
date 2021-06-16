// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package archive

import (
	"fmt"
	zip2 "github.com/ainsleyclark/updater/internal/archive/zip"
	"github.com/ainsleyclark/updater/internal/fileio"
	"strings"
)

type Archive interface {
	Copy(dest string) error
}

func New(src string) (Archive, error) {
	if !fileio.Exists(src) {
		return nil, fmt.Errorf("file not found with the path: %s", src)
	}
	if strings.HasSuffix(src, ".zip") {
		return &zip2.Zip{Path: src}, nil
	}
	return nil, fmt.Errorf("no archiver found")
}
