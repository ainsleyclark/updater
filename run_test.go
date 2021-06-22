// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const (
	v001 = "UPDATE my_table SET name = 'tom' WHERE id = 1"
	v002 = "UPDATE my_table SET name = 'nick' WHERE id = 2"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

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
				&Migration{Version: "v0.0.1", SQL: strings.NewReader(v001), Stage: Major},
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
				&Migration{Version: "v0.0.1", SQL: strings.NewReader(v001), Stage: Major},
			},
			func(m sqlmock.Sqlmock) {
				m.ExpectBegin().
					WillReturnError(fmt.Errorf("error"))
			},
			true,
			"error",
			DatabaseError,
		},
		"Exec Error": {
			migrationRegistry{
				&Migration{Version: "v0.0.1", SQL: strings.NewReader(v001), Stage: Major},
			},
			func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(v001).
					WillReturnError(fmt.Errorf("error"))
				m.ExpectRollback()
			},
			true,
			"error",
			DatabaseError,
		},
		"RollBack Error": {
			migrationRegistry{
				&Migration{Version: "v0.0.1", SQL: strings.NewReader(v001), Stage: Major},
			},
			func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(v001).
					WillReturnError(fmt.Errorf("error"))
				m.ExpectRollback().
					WillReturnError(fmt.Errorf("error"))
			},
			true,
			"error",
			DatabaseError,
		},
		"Commit Error": {
			migrationRegistry{
				&Migration{Version: "v0.0.1", SQL: strings.NewReader(v001), Stage: Major},
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
				&Migration{Version: "v0.0.0", SQL: strings.NewReader(v001), Stage: Major},
			},
			nil,
			false,
			"error",
			Updated,
		},
		"Bad SQL": {
			migrationRegistry{
				&Migration{Version: "v0.0.1", SQL: errReader(1), Stage: Major},
			},
			nil,
			false,
			"error",
			Unknown,
		},
		"With Callback": {
			migrationRegistry{
				&Migration{Version: "v0.0.1", SQL: strings.NewReader(v001), Stage: Major, CallBackUp: func() error {
					return nil
				}, CallBackDown: func() error {
					return nil
				}},
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
		"With Callback Up Error": {
			migrationRegistry{
				&Migration{Version: "v0.0.1", SQL: strings.NewReader(v001), Stage: Major, CallBackUp: func() error {
					return fmt.Errorf("error")
				}, CallBackDown: func() error {
					return nil
				}},
			},
			func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(v001).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectRollback()
			},
			true,
			"error",
			CallBackError,
		},
		"With Callback Down Error": {
			migrationRegistry{
				&Migration{Version: "v0.0.1", SQL: strings.NewReader(v001), Stage: Major, CallBackUp: func() error {
					return nil
				}, CallBackDown: func() error {
					return fmt.Errorf("callback error")
				}},
				&Migration{Version: "v0.0.2", SQL: strings.NewReader(v002), Stage: Major, CallBackUp: func() error {
					return fmt.Errorf("error")
				}, CallBackDown: func() error {
					return nil
				}},
			},
			func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectExec(v001).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectExec(v002).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectRollback()
			},
			true,
			"callback error",
			CallBackError,
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
					Version:       "0.0.0",
					RepositoryURL: "https://github.com/ainsleyclark/verbis",
					hasDB: test.db,
				},
				pkg:     nil,
				version: version.Must(version.NewVersion("0.0.0")),
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
