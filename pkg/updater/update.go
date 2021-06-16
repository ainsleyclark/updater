// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"errors"
	"github.com/ainsleyclark/updater/internal/archive"
	"github.com/ainsleyclark/updater/internal/fileio"
	"github.com/ainsleyclark/updater/internal/patcher"
	"github.com/ainsleyclark/updater/pkg/github"
	"github.com/mattn/go-zglob"
	"io/fs"
	"os"
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
	mover patcher.Mover
}

const (
	DefaultBackupExtension = ".bak"
)

var (
	ErrLatestVersion = errors.New("at latest version")
)

func (u *Updater) Update() (err error) {
	update, err := u.CanUpdate()
	if err != nil || !update {
		return
	}

	if u.BackupExtension == "" {
		u.BackupExtension = DefaultBackupExtension
	}

	tempZip, err := u.Github.Download()
	if err != nil {
		return err
	}
	defer u.Github.Close()

	a, err := archive.New(tempZip)
	if err != nil {
		return err
	}

	tmpArchive, err := fileio.TempDirectory()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpArchive)

	err = a.Copy(tmpArchive)
	if err != nil {
		return err
	}

	u.mover = patcher.New()

	err = u.updateExecutable(tmpArchive)
	if err != nil {
		return err
	}

	//err = u.updateFiles(tmpArchive)
	//if err != nil {
	//	return err
	//}

	return u.mover.Apply()
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

//nolint
func (u *Updater) updateExecutable(archiveDir string) error {
	exec, err := fileio.Executable()
	if err != nil {
		return err
	}

	err = filepath.WalkDir(archiveDir, func(path string, d fs.DirEntry, err error) error {
		relativePath := strings.ReplaceAll(path, archiveDir+string(os.PathSeparator), "")
		if relativePath == u.RemoteExecutablePath {
			p := patcher.Patcher{
				DestinationPath: exec,
				BackupPath:      exec + u.BackupExtension,
			}
			p.AddFile(&patcher.File{
				SourcePath:      path,
				DestinationPath: exec,
				Mode:            0755,
			})
			u.mover.AddPatcher(p)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *Updater) updateFiles(archiveDir string) error {
	for _, f := range u.Files {
		src := filepath.Join(archiveDir, f.RemotePath)

		// This should be a validate function of files
		matches, err := zglob.Glob(src)
		if err != nil {
			return err
		}

		p := patcher.Patcher{
			DestinationPath: f.LocalPath,
			BackupPath:      f.LocalPath + u.BackupExtension,
		}

		for _, v := range matches {
			cleanedFile := strings.ReplaceAll(v, archiveDir+string(os.PathSeparator), "")
			path, err := fileio.SplitPaths(f.RemotePath, cleanedFile)
			if err != nil {
				return err
			}
			stat, err := os.Stat(v)
			if err != nil {
				return err
			}
			p.AddFile(&patcher.File{
				SourcePath:      v,
				DestinationPath: filepath.Join(f.LocalPath, path),
				Mode:            stat.Mode(),
			})
		}

		u.mover.AddPatcher(p)
	}
	
	return nil
}
