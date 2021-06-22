// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"errors"
	"github.com/hashicorp/go-version"
	"io"
	"sort"
)

// Migration represents a singular migration for a single
// version.
type Migration struct {
	// The main version of the migration such as "v0.0.1"
	Version string
	// The migration file, byte value of the SQL migration.
	SQL io.Reader
	// CallBackUp is a function called when the migration
	// is going up, this can be useful when manipulating
	// files and directories for the current version.
	CallBackUp CallBackFn
	// CallBackUp is a function called when the migration
	// is going down, this is only called if an update
	// failed. And must be provided if CallBackUp is
	// defined.
	CallBackDown CallBackFn
	// Stage defines the release stage of the migration such as
	// Major, Minor or Patch,
	Stage Stage
}

// CallBackFn is the function type when migrations are
// running up or down.
type CallBackFn func() error

// migrationRegistry contains a slice of pointers to each
// migration.
type migrationRegistry []*Migration

// migrations is the in memory registry store of
// migrations.
var migrations = make(migrationRegistry, 0)

var (
	// ErrCallBackMismatch is returned by AddMigration when
	// there has been a mismatch in the amount of callbacks
	// passed. Each migration should have two callbacks,
	// one up and one down, or none at all.
	ErrCallBackMismatch = errors.New("both CallBackUp and CallBackDown must be set")
)

// GetMigration Retrieves a migration from the registry by
// looking up the version. An error will be returned on
// failed lookup.
func GetMigration(version string) (*Migration, error) { //nolint
	for _, m := range migrations {
		if m.Version == version {
			return m, nil
		}
	}
	return nil, errors.New("no migration found with the version: " + version)
}

// AddMigration adds a migration to the update registry
// which will be called when Update() is run. The
// version and Stage must be attached to the
// migration.
func AddMigration(m *Migration) error {
	if m.Version == "" {
		return errors.New("no version provided for update")
	}

	if m.Stage == "" {
		return errors.New("no stage set")
	}

	if m.CallBackUp != nil && m.CallBackDown == nil {
		return ErrCallBackMismatch
	}

	if m.CallBackUp == nil && m.CallBackDown != nil {
		return ErrCallBackMismatch
	}

	_ = m.toSemVer() // Check to see if the version is valid

	migrations = append(migrations, m)

	return nil
}

// toSemVer parses the migration to version.Version, and
// panics if the version is malformed.
func (m *Migration) toSemVer() *version.Version {
	semver, err := version.NewVersion(m.Version)
	if err != nil {
		panic(err.Error())
	}
	return semver
}

// hasCallBack returns true if CallBackUp and CallBackDown
// are both defined.
func (m *Migration) hasCallBack() bool {
	return m.CallBackUp != nil && m.CallBackDown != nil
}

// Sort migrationRegistry is a type that implements the
// sort.Interface interface so that versions can be
// sorted.
func (m migrationRegistry) Sort() {
	sort.Sort(m)
}

func (m migrationRegistry) Len() int {
	return len(m)
}

func (m migrationRegistry) Less(i, j int) bool {
	return m[i].toSemVer().LessThan(m[j].toSemVer())
}

func (m migrationRegistry) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
