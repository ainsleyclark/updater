// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package patcher

import (
	"io/ioutil"
	"os"
)

type Patcher struct {
	SourcePath      string
	DestinationPath string
	BackupPath      string
	Mode            os.FileMode
	IsDirectory     bool
}


func (p *Patcher) Apply() error {
	content, err := ioutil.ReadFile(p.SourcePath)
	if err != nil {
		return err
	}

	//err = p.Backup()
	//if err != nil {
	//	return err
	//}


	err = ioutil.WriteFile(p.DestinationPath, content, p.Mode)
	if err != nil {
		return p.Rollback()
	}

	return nil
}

// Backup renames a directory or file to the new path.
func (p *Patcher) Backup() error {
	return os.Rename(p.DestinationPath, p.BackupPath)
}

func (p *Patcher) Rollback() error {
	return os.Rename(p.BackupPath, p.DestinationPath)
}

