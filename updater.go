// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"github.com/hashicorp/go-version"
	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
)

// |||||||||||||||||||||||||||||||||||||||||||||||||||||||
type Patcher interface {
	Update(archive string) (Status, error)
	HasUpdate() (bool, error)
	LatestVersion() (string, error)
}

// |||||||||||||||||||||||||||||||||||||||||||||||||||||||
type Updater struct {
	opts    Options
	pkg     *updater.Updater
	version *version.Version
}

// New returns a new Updater with the options passed. If
// validation failed on the options an error will be
// returned.
func New(opts Options) (*Updater, error) {
	err := opts.Validate()
	if err != nil {
		return nil, err
	}

	ver, err := version.NewVersion(opts.Version)
	if err != nil {
		return nil, err
	}

	u := &Updater{
		opts: opts,
		pkg: &updater.Updater{
			Provider: &provider.Github{
				RepositoryURL: opts.RepositoryURL,
			},
			Version: opts.Version,
		},
		version: ver,
	}

	return u, nil
}

// HasUpdate determines if there is an update for the
// program. Returns a error if there are no releases
// or tags for the repo.
func (u *Updater) HasUpdate() (bool, error) {
	return u.pkg.CanUpdate()
}

// LatestVersion retrieves the most up to date version of
// the program. Returns a error if there are no releases
//// or tags for the repo.
func (u *Updater) LatestVersion() (string, error) {
	return u.pkg.GetLatestVersion()
}

// Update takes in the archive name of the zip file or
// folder to download and proceeds to update the
// executable and migrates any database queries
// or callbacks. If there was an error in any
// of the processes, the package will
// rollback to the previous state.
func (u *Updater) Update(archive string) (Status, error) {
	u.pkg.Provider = &provider.Github{
		RepositoryURL: u.opts.RepositoryURL,
		ArchiveName:   archive,
	}

	update, err := u.pkg.Update()
	status := getExecStatus(update)

	if err != nil {
		return status, err
	}

	if u.opts.Verify {
		err = u.verifyInstallation()
		if err != nil {
			return ExecutableError, err
		}
	}

	// Run any migrations

	return status, nil
}
