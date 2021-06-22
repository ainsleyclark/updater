// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

// Options define the core arguments parsed to the migrator.
type Options struct {
	// The URL of the GitHub Repository to obtain the
	// executable from.
	RepositoryURL string
	// The currently running version.
	Version string
	// If set to true, updates will be verified by checking the
	// newly downloaded executable version number using the
	// -version flag.
	Verify bool
	// SQL database to apply migrations, migrations will not
	// be run if sql.DB is nil.
	DB *sql.DB
	// Determines if the database is set.
	hasDB bool
}

var (
	// ErrRepositoryURL is the error returned by Validate when
	// a malformed repository is used.
	ErrRepositoryURL = errors.New("error checking repo url")
)

// Validate check's to see if the options are valid before
// returning a new migrator.
func (o *Options) Validate() error {
	if o.RepositoryURL == "" {
		return errors.New("no repo url provided")
	}

	if o.Version == "" {
		return errors.New("no version provided")
	}

	resp, err := http.Get(o.RepositoryURL)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrRepositoryURL, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return ErrRepositoryURL
	}

	if o.DB != nil {
		err := o.DB.Ping()
		if err != nil {
			return err
		}
		o.hasDB = true
	}

	return nil
}
