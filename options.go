// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import "errors"

type Options struct {
	RepositoryURL string
	ArchiveName   string
	Version       string
}

func (o *Options) Validate() error {
	if o.RepositoryURL == "" {
		return errors.New("no repo url provided")
	}
	if o.ArchiveName == "" {
		return errors.New("no archive name provided")
	}
	if o.Version == "" {
		return errors.New("no version provided")
	}
	return nil
	// TODO: check if its an ok http response
}
