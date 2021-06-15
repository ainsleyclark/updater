// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"errors"
	"fmt"
	"github.com/ainsleyclark/updater/internal/fileio"
	"github.com/ainsleyclark/updater/internal/patcher"
	"github.com/ainsleyclark/updater/pkg/archive"
	"github.com/ainsleyclark/updater/pkg/github"
	"github.com/gookit/color"
	"github.com/mattn/go-zglob"
	"io/fs"
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


func (u *Updater) Update() (err error) {
	update, err := u.CanUpdate()
	if err != nil || !update {
		return
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

	u.test(tmpArchive)
	//
	//err = u.walk(tmpArchive)
	//if err != nil {
	//	return err
	//}

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

func (u *Updater) walk(src string) error {
	err := filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
	 	if entry.IsDir() {
			return nil
		}

		relativePath := strings.ReplaceAll(path, src + string(os.PathSeparator), "")

		//// Check for executable and update
		//if relativePath == u.RemoteExecutablePath {
		//	err := u.updateExecutable(path)
		//	if err != nil {
		//		return err
		//	}
		//}

		// Check for files and folders and update
		err = u.updateFiles(src, relativePath, entry)
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

func (u *Updater) updateExecutable(path string) error {
	exec, err := fileio.Executable()
	if err != nil {
		return err
	}

	// Need a health check
	if u.BackupExtension == "" {
		u.BackupExtension = ".bak"
	}

	p := patcher.Patcher{
		SourcePath:      path,
		DestinationPath: exec,
		BackupPath:     exec + u.BackupExtension,
		Mode:            0755,
	}

	return p.Apply()
}

func (u *Updater) updateFiles(tmpDir string, p string, entry fs.DirEntry) error {
	var matches []string
	for _, f := range u.Files {
		match, err := zglob.Match(f.RemotePath, p)
		if err != nil {
			return err
		}

		if !match {
			continue
		}

		// TODO:
		// Account for files only e.g index.html

		relative, err := filepath.Rel(f.RemotePath, p)
		if err != nil {
			return err
		}

		dest := path.Join(f.LocalPath, strings.ReplaceAll(relative, "../", ""))

		matches = append(matches, dest)
	}




	return nil
}


func (u *Updater) test(archiveDir string) {
	//var result = make(map[string][]string)
	for _, f := range u.Files {
		matches, err := zglob.Glob(filepath.Join(archiveDir,f.RemotePath))
		if err != nil {
			return
		}

		for _, v := range matches {
			cleanedFile := strings.ReplaceAll(v, archiveDir + string(os.PathSeparator), "")

			path, err := u.relativePath(f.RemotePath, cleanedFile)
			if err != nil {
				return
			}

			color.Yellow.Println("Relative: ", path)
			color.Red.Println("Source: ", v)
			color.Blue.Println("Destination: ", filepath.Join(f.LocalPath, path))
			fmt.Println("-------------")
		}
	}

}

func (u *Updater) relativePath(base, target string) (string, error) {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return "", err
	}
	if rel == "." {
		return "", nil
	}
	return strings.ReplaceAll(rel, "../", ""), nil
}

