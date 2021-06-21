// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	sm "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMigration(t *testing.T) {
	tt := map[string]struct {
		input     string
		migration Migration
		want      interface{}
	}{
		"Found": {
			"v0.0.1",
			Migration{Version: "v0.0.1"},
			Migration{Version: "v0.0.1"},
		},
		"Not Found": {
			"v0.0.1",
			Migration{Version: "v0.0.2"},
			"no migration found with the version",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			migrations = []*Migration{&test.migration}
			defer func() {
				migrations = make(migrationRegistry, 0)
			}()
			got, err := GetMigration(test.input)
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}
			assert.Equal(t, test.want, *got)
		})
	}
}

func TestAddMigration(t *testing.T) {
	tt := map[string]struct {
		input Migration
		want  interface{}
	}{
		"Success": {
			Migration{Version: "v0.0.1", Stage: Minor, MigrationPath: "v0.0.1.sql"},
			nil,
		},
		"No version": {
			Migration{Version: ""},
			"no version provided for update",
		},
		"No Stage": {
			Migration{Version: "v0.0.1"},
			"no stage set",
		},
		"Bad version": {
			Migration{Version: "v1.3.3.3", MigrationPath: "test", Stage: Minor},
			"invalid version",
		},
		"No CallBackUp": {
			Migration{Version: "v0.0.1", MigrationPath: "test", Stage: Minor, CallBackDown: func() error {
				return nil
			}},
			ErrCallBackMismatch.Error(),
		},
		"No CallBackDown": {
			Migration{Version: "v0.0.1", MigrationPath: "test", Stage: Minor, CallBackUp: func() error {
				return nil
			}},
			ErrCallBackMismatch.Error(),
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			defer func() {
				migrations = make(migrationRegistry, 0)
			}()
			err := AddMigration(&test.input)
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}
			assert.Equal(t, test.input, *migrations[0])
		})
	}
}

func TestMigrationRegistry_Sort(t *testing.T) {
	tt := map[string]struct {
		input migrationRegistry
		want  migrationRegistry
	}{
		"Success": {
			migrationRegistry{&Migration{Version: "v3.0.0"}, &Migration{Version: "v1.0.0"}, &Migration{Version: "v2.0.0"}},
			migrationRegistry{&Migration{Version: "v1.0.0"}, &Migration{Version: "v2.0.0"}, &Migration{Version: "v3.0.0"}},
		},
		"Already Sorted": {
			migrationRegistry{&Migration{Version: "v1.0.0"}, &Migration{Version: "v2.0.0"}, &Migration{Version: "v3.0.0"}},
			migrationRegistry{&Migration{Version: "v1.0.0"}, &Migration{Version: "v2.0.0"}, &Migration{Version: "v3.0.0"}},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			test.input.Sort()
			for i, v := range test.input {
				assert.Equal(t, test.want[i].Version, v.Version)
			}
		})
	}
}

func TestMigration_ToSemVer(t *testing.T) {
	tt := map[string]struct {
		input  Migration
		panics bool
		want   interface{}
	}{
		"Success": {
			Migration{Version: "v0.0.1"},
			false,
			sm.Must(sm.NewVersion("v0.0.1")),
		},
		"Error": {
			Migration{Version: "wrong"},
			true,
			nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			if test.panics {
				assert.Panics(t, func() {
					got := test.input.toSemVer()
					assert.Equal(t, test.want, got)
				})
				return
			}
			assert.Equal(t, test.want, test.input.toSemVer())
		})
	}
}

func TestMigration_HasCallBack(t *testing.T) {
	tt := map[string]struct {
		input Migration
		want  bool
	}{
		"Has CallBack": {
			Migration{CallBackUp: func() error {
				return nil
			}, CallBackDown: func() error {
				return nil
			}},
			true,
		},
		"No Callback": {
			Migration{},
			false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.input.hasCallBack()
			assert.Equal(t, test.want, got)
		})
	}
}
