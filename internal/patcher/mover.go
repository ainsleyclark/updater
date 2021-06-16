// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package patcher

import "os"

type Mover struct {
	p         []Patcher
	completed []Patcher
}

func New() Mover {
	return Mover{
		p:         make([]Patcher, 0),
		completed: make([]Patcher, 0),
	}
}

func (m *Mover) Apply() error {
	for _, v := range m.p {
		err := v.apply()
		if err != nil {
			return m.RollBack()
		}
		m.completed = append(m.completed, v)
	}
	return m.CleanUp()
}

func (m *Mover) AddPatcher(p Patcher) {
	if len(m.p) == 0 {
		m.p = make([]Patcher, 0)
	}
	m.p = append(m.p, p)
}

func (m *Mover) RollBack() error {
	for _, p := range m.completed {
		err := os.RemoveAll(p.DestinationPath)
		if err != nil {
			return err
		}
		if p.backupPossible {
			continue
		}
		err = os.Rename(p.BackupPath, p.DestinationPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Mover) CleanUp() error {
	for _, p := range m.completed {
		if !p.backupPossible && !p.IsExec {
			continue
		}
		err := os.RemoveAll(p.BackupPath)
		if err != nil {
			return err
		}
	}
	return nil
}
