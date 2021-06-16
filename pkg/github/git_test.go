// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"github.com/ainsleyclark/updater/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRepo_LatestVersion(t *testing.T) {
	tt := map[string]struct {
		apiUrl string
		input  string
		want   interface{}
	}{
		"Success": {
			apiUrl,
			tests.ReleaseRepo,
			nil,
		},
		"Bad Input URL": {
			apiUrl,
			"wrong",
			"invalid github url",
		},
		"Bad API URL": {
			"wrong",
			tests.ReleaseRepo,
			"unsupported protocol scheme",
		},
		"No Tags": {
			apiUrl,
			"github.com/ainsleyclark/html-boilerplate",
			"repo has no tags",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			orig := apiUrl
			defer func() {
				apiUrl = orig
			}()
			apiUrl = test.apiUrl

			r := Repo{RepositoryURL: test.input}
			got, err := r.LatestVersion()
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}
			assert.NotEmpty(t, got)
		})
	}
}

func TestRepo_Close(t *testing.T) {
	tt := map[string]struct {
		input string
		want  interface{}
	}{
		"Success": {
			t.TempDir(),
			nil,
		},
		"Remove Error": {
			".",
			"error removing temp directory",
		},
		"No Temp Dir": {
			"",
			nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			r := Repo{tempDir: test.input}
			got := r.Close()
			if got != nil {
				require.Contains(t, got.Error(), test.want)
				return
			}
			require.Equal(t, test.want, got)
		})
	}
}

func TestRepo_GetInfo(t *testing.T) {
	tt := map[string]struct {
		input string
		want  interface{}
	}{
		"Success": {
			"https://github.com/tom/repo",
			information{Owner: "tom", Name: "repo"},
		},
		"Invalid URL": {
			"wrong",
			"invalid github url",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			r := Repo{RepositoryURL: test.input}
			err := r.getInfo()
			if err != nil {
				require.Contains(t, err.Error(), test.want)
				return
			}
			require.Equal(t, test.want, r.info)
		})
	}
}

func TestRepo_GetDownloadURL(t *testing.T) {
	tt := map[string]struct {
		info information
		tag  string
		name string
		want interface{}
	}{
		"Simple": {
			information{Owner: "tom", Name: "repo"},
			"0.0.1",
			"archive.zip",
			"https://github.com/tom/repo/releases/download/0.0.1/archive.zip",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			r := Repo{info: test.info}
			got := r.getDownloadURL(test.tag, test.name)
			require.Equal(t, test.want, got)
		})
	}
}

func TestRepo_GetTags(t *testing.T) {
	tt := map[string]struct {
		handler http.HandlerFunc
		want    interface{}
	}{
		"Success": {
			func(w http.ResponseWriter, r *http.Request) {
				bytes, err := json.Marshal([]tag{{Name: "0.0.1"}})
				require.NoError(t, err)
				_, err = w.Write(bytes)
				require.NoError(t, err)
			},
			[]tag{{Name: "0.0.1"}},
		},
		"Mismatch Status": {
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte("error"))
				require.NoError(t, err)
			},
			"error",
		},
		"Unmarshal Error": {
			func(w http.ResponseWriter, r *http.Request) {},
			"unexpected end of JSON input",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(test.handler)
			defer ts.Close()

			orig := apiUrl
			defer func() {
				apiUrl = orig
			}()
			apiUrl = ts.URL

			r := Repo{}
			got, err := r.getTags()
			if err != nil {
				require.Contains(t, err.Error(), test.want)
				return
			}

			require.Equal(t, test.want, got)
		})
	}
}

func TestRepo_GetTagsErorr(t *testing.T) {
	orig := apiUrl
	defer func() {
		apiUrl = orig
	}()
	apiUrl = "wrong"
	r := Repo{}
	_, err := r.getTags()
	require.Error(t, err)
}
