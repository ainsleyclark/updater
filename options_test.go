// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOptions_Validate(t *testing.T) {
	tt := map[string]struct {
		input   Options
		handler func(w http.ResponseWriter, r *http.Request)
		want    interface{}
	}{
		"Success": {
			Options{RepositoryURL: "/migrator", Version: "0.0.1"},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			nil,
		},
		"No Repo": {
			Options{Version: "0.0.1"},
			nil,
			"no repo url provided",
		},
		"No Version": {
			Options{RepositoryURL: "url"},
			nil,
			"no version provided",
		},
		"Bad URL": {
			Options{RepositoryURL: "https://", Version: "0.0.1"},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			"no such host",
		},
		"Invalid Status Code": {
			Options{RepositoryURL: "/migrator", Version: "0.0.1"},
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			ErrRepositoryURL.Error(),
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			if test.handler != nil {
				ts := httptest.NewServer(http.HandlerFunc(test.handler))
				defer ts.Close()
				test.input.RepositoryURL = ts.URL + test.input.RepositoryURL
			}
			err := test.input.Validate()
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
			}
		})
	}
}
