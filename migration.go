// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"errors"
	"github.com/hashicorp/go-version"
	"sort"
)

type Migration struct {
	Version       string
	MajorVersion  int
	MigrationPath string
	CallBackUp    CallBackFn
	CallBackDown  CallBackFn
	Stage         Stage
}

type CallBackFn func() error

type migrationRegistry []*Migration

var migrations = make(migrationRegistry, 0)

var (
	ErrCallBackMismatch = errors.New("both CallBackUp and CallBackDown must be set")
)

// AddMigration add's a migration to the update registry which will be called when Update() is run.
//
//. the

func AddMigration(m *Migration) error {
	if m.Version == "" {
		return errors.New("no version provided for update")
	}

	if m.Stage == "" {
		return errors.New("no stage set")
	}

	if m.MigrationPath == "" {
		return errors.New("no migration path set")
	}

	if m.CallBackUp != nil && m.CallBackDown == nil {
		return ErrCallBackMismatch
	}

	if m.CallBackUp == nil && m.CallBackDown != nil {
		return ErrCallBackMismatch
	}

	semVer := m.toSemVer()
	seg := semVer.Segments()

	if len(seg) != 3 { //nolint
		return errors.New("invalid version: " + m.Version)
	}

	m.MajorVersion = seg[0]

	migrations = append(migrations, m)

	return nil
}

func (m *Migration) toSemVer() *version.Version {
	semver, err := version.NewVersion(m.Version)
	if err != nil {
		panic(err.Error())
	}
	return semver
}

// Sort migrationRegistry is a type that implements the sort.Interface
// interface so that versions can be sorted.
func (r migrationRegistry) Sort() {
	sort.Sort(r)
}

func (r migrationRegistry) Len() int {
	return len(r)
}

func (r migrationRegistry) Less(i, j int) bool {
	return r[i].toSemVer().LessThan(r[j].toSemVer())
}

func (r migrationRegistry) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (m *Migration) hasCallBack() bool {
	return m.CallBackUp != nil && m.CallBackDown != nil
}
