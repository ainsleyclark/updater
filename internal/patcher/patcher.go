// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package patcher

import (
	"errors"
	"github.com/ainsleyclark/updater/internal/fileio"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Patcher struct {
	DestinationPath string
	BackupPath      string
	IsExec          bool
	files           []*File
	backupPossible  bool
}

type File struct {
	SourcePath      string
	DestinationPath string
	Mode            os.FileMode
}

func (p *Patcher) apply() error {
	if fileio.Exists(p.DestinationPath) && !p.IsExec {
		p.backupPossible = true
	}

	if fileio.Exists(p.BackupPath) {
		return errors.New("backup path already exists: " + p.BackupPath)
	}

	err := p.backup()
	if err != nil {
		return err
	}

	for _, f := range p.files {
		err := os.MkdirAll(filepath.Dir(f.DestinationPath), os.ModePerm)
		if err != nil {
			return err
		}
		if f.Mode.IsDir() {
			continue
		}
		content, err := ioutil.ReadFile(f.SourcePath)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(f.DestinationPath, content, f.Mode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Patcher) AddFile(f *File) {
	if len(p.files) == 0 {
		p.files = make([]*File, 0)
	}
	if f == nil {
		return
	}
	p.files = append(p.files, f)
}

// backup renames a directory or file to the new path.
func (p *Patcher) backup() error {
	if p.backupPossible {
		return os.Rename(p.DestinationPath, p.BackupPath)
	}
	return nil
}
