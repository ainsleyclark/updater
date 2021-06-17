// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
)

type Patcher interface {
	Update() (Status, error)
	HasUpdate() (bool, error)
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
				ArchiveName:   opts.ArchiveName,
			},
			Version: opts.Version,
		},
	}, nil
}

func (u *Updater) HasUpdate() (bool, error) {
	return u.pkg.CanUpdate()
}

func (u *Updater) Update() (Status, error) {
	update, err := u.pkg.Update()
	status := getExecStatus(update)
	if err != nil {
		return status, err
	}
	return status, nil
}