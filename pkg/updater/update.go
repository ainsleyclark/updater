// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"errors"
	"github.com/ainsleyclark/updater/pkg/fileio"
	"github.com/ainsleyclark/updater/pkg/github"
	"github.com/ainsleyclark/updater/pkg/patcher"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// |||||||||||||||||||||||||||||||||||||||||||||||||||||||

type Updater struct {
	Github               github.Repo
	Files                Files
	Version              string
	BackupExtension      string
	RemoteExecutablePath string
}

var (
	ErrLatestVersion = errors.New("at latest version")
)

type Options struct {
}

func New() {

}

func (u *Updater) Update() (err error) {
	update, err := u.CanUpdate()
	if err != nil || !update {
		return
	}

	err = u.Github.Open()
	if err != nil {
		return
	}
	defer u.Github.Close()

	err = u.walk()
	if err != nil {
		return
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

// get the currently running executable name
//

func (u *Updater) walk() error {
	err := u.Github.Walk(func(info *github.FileInfo) error {
		// Update executable
		if info.Mode.IsRegular() && info.Path == u.RemoteExecutablePath {
			err := u.updateExecutable(info)
			if err != nil {
				return err
			}
		}
		// Update any files
		err := u.updateFiles(info)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *Updater) updateExecutable(info *github.FileInfo) error {
	tmp, err := fileio.TempDirectory()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	paths, err := fileio.GetPaths()
	if err != nil {
		return err
	}

	tmpExec := filepath.Join(tmp, paths.Base)
	err = u.Github.Copy(info.Path, tmpExec)
	if err != nil {
		return err
	}

	// Need a health check
	if u.BackupExtension == "" {
		u.BackupExtension = ".bak"
	}

	p := patcher.Patcher{
		SourcePath:      tmpExec,
		DestinationPath: paths.ExecutableName,
		BackupPath:      paths.ExecutableName + u.BackupExtension,
		Mode:            0755,
	}

	return p.Apply()
}

func (u *Updater) updateFiles(info *github.FileInfo) error {
	for _, f := range u.Files {
		match, err := path.Match(f.RemotePath, info.Path)
		if err != nil {
			return err
		}

		if !match {
			continue
		}

		relative, err := filepath.Rel(f.RemotePath, info.Path)
		if err != nil {
			return err
		}

		dest := path.Join(f.LocalPath, strings.ReplaceAll(relative, "../", ""))
		err = os.MkdirAll(path.Dir(dest), os.ModePerm)
		if err != nil {
			return err
		}

		err = u.Github.Copy(info.Path, dest)
		if err != nil {
			return err
		}
	}
	return nil
}
