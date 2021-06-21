// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
)

type Patcher interface {
	Update(archive string) (Status, error)
	HasUpdate() (bool, error)
	LatestVersion() (string, error)
}

type Updater struct {
	opts *Options
	pkg  *updater.Updater
}

func New(opts *Options) (*Updater, error) {
	err := opts.Validate()
	if err != nil {
		return nil, err
	}
	return &Updater{
		opts: opts,
		pkg: &updater.Updater{
			Provider: &provider.Github{
				RepositoryURL: opts.RepositoryURL,
				ArchiveName:   "",
			},
		},
	}, nil
}

func (u *Updater) HasUpdate() (bool, error) {
	return u.pkg.CanUpdate()
}

func (u *Updater) LatestVersion() (string, error) {
	return u.pkg.GetLatestVersion()
}

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
