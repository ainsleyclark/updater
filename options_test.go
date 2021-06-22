// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOptions_Validate(t *testing.T) {
	tt := map[string]struct {
		input   Options
		handler func(w http.ResponseWriter, r *http.Request)
		db func(mock sqlmock.Sqlmock)
		want    interface{}
	}{
		"Success": {
			Options{RepositoryURL: "/migrator", Version: "0.0.1"},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			nil,
			nil,
		},
		"No Repo": {
			Options{Version: "0.0.1"},
			nil,
			nil,
			"no repo url provided",
		},
		"No version": {
			Options{RepositoryURL: "url"},
			nil,
			nil,
			"no version provided",
		},
		"Bad URL": {
			Options{RepositoryURL: "https://", Version: "0.0.1"},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			nil,
			"no such host",
		},
		"Invalid Status Code": {
			Options{RepositoryURL: "/migrator", Version: "0.0.1"},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			nil,
			ErrRepositoryURL.Error(),
		},
		"With DB": {
			Options{RepositoryURL: "/migrator", Version: "0.0.1"},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			func(mock sqlmock.Sqlmock) {
				mock.ExpectPing()
			},
			ErrRepositoryURL.Error(),
		},
		 "Ping Error": {
			Options{RepositoryURL: "/migrator", Version: "0.0.1"},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			func(mock sqlmock.Sqlmock) {
				mock.ExpectPing().
					WillReturnError(fmt.Errorf("ping error"))
			},
			"ping error",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			if test.handler != nil {
				ts := httptest.NewServer(http.HandlerFunc(test.handler))
				defer ts.Close()
				test.input.RepositoryURL = ts.URL + test.input.RepositoryURL
			}

			var mock sqlmock.Sqlmock

			if test.db != nil {
				db, m, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
				assert.NoError(t, err)
				mock = m
				test.db(mock)
				test.input.DB = db
			}

			err := test.input.Validate()
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
			}

			if test.db != nil {
				err = mock.ExpectationsWereMet()
				if err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}
}
