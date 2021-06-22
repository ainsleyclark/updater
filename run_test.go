// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const (
	v001 = "UPDATE my_table SET name = 'tom' WHERE id = 1"
)

func TestUpdater_Run(t *testing.T) {
	tt := map[string]struct {
		input migrationRegistry
		mock  func(m sqlmock.Sqlmock)
		db bool
		want  interface{}
		code Status
	}{
		"Simple": {
			migrationRegistry{
				&Migration{Version: "v0.0.1", Migration: strings.NewReader(v001), Stage: Major},
			},
			func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(v001).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
			true,
			nil,
			Updated,
		},
		"Begin Error": {
			migrationRegistry{
				&Migration{Version: "v0.0.1", Migration: strings.NewReader(v001), Stage: Major},
			},
			func(m sqlmock.Sqlmock) {
				m.ExpectBegin().
					WillReturnError(fmt.Errorf("error"))
			},
			true,
			"error",
			DatabaseError,
		},
		"Commit Error": {
			migrationRegistry{
				&Migration{Version: "v0.0.1", Migration: strings.NewReader(v001), Stage: Major},
			},
			func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(v001).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit().
					WillReturnError(fmt.Errorf("error"))
			},
			true,
			"error",
			DatabaseError,
		},

		"No Run": {
			migrationRegistry{
				&Migration{Version: "v0.0.0", Migration: strings.NewReader(v001), Stage: Major},
			},
			nil,
			false,
			"error",
			Updated,
		},
		"Bad Migration": {
			migrationRegistry{
				&Migration{Version: "v0.0.0", Migration: nil, Stage: Major},
			},
			nil,
			false,
			"error",
			Updated,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			if test.mock != nil {
				test.mock(mock)
			}

			defer func() {
				migrations = make(migrationRegistry, 0)
				db.Close()
			}()

			u := Updater{
				opts:    Options{
					DB:            db,
					Version:       "0.0.1",
					RepositoryURL: "https://github.com/ainsleyclark/verbis",
					hasDB: test.db,
				},
				pkg:     nil,
				version: version.Must(version.NewVersion("0.0.1")),
			}

			migrations = test.input
			assert.NoError(t, err)

			code, err := u.runMigrations()
			assert.Equal(t, test.code, code)
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}

			if test.db {
				err = mock.ExpectationsWereMet()
				if err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}
}
