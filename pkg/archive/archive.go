// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package archive

type Archiver interface {
	Open() error
	Close() error
	Unzip(dest string) error
}