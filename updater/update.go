// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"errors"
	"github.com/ainsleyclark/updater/github"
	"os"
)

type Updater struct {
	Github          github.Repo
	Files           Files
	ExecutableName  string
	Version         string
	BackupExtension string
}

var (
	ErrLatestVersion = errors.New("at latest version")
)

func (u *Updater) Update() error {
	update, err := u.CanUpdate()
	if err != nil || !update {
		return err
	}

	err = u.Files.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (u *Updater) LatestVersion() (string, error) {
	version, err := u.Github.LatestVersion()
	if err != nil {
		return "", err
	}
	return version, nil
}

func (u *Updater) CanUpdate() (bool, error) {
	version, err := u.Github.LatestVersion()
	if err != nil {
		return false, err
	}
	if u.Version == version {
		return false, ErrLatestVersion
	}
	return true, nil
}


type Files []string

func (f Files) Validate() error {
	for _, file := range f {
		_, err := os.Stat(file)
		if !os.IsNotExist(err) {
			return errors.New("no file or directory exists with the path: " + file)
		}
	}
	return nil
}

