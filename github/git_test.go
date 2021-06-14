// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"github.com/ainsleyclark/updater/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepo_LatestVersion(t *testing.T) {
	tt := map[string]struct {
		input string
		want interface{}
	}{
		"Success": {
			tests.ReleaseRepo,
			nil,
		},
		"Bad URL": {
			"wrong",
			"invalid github URL",
		},
		"No Tags": {
			"github.com/ainsleyclark/html-boilerplate",
			"repo has no tags",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
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

func TestRepo_Walk(t *testing.T) {
	tt := map[string]struct {
		open bool
		fn func(info *FileInfo) error
		want interface{}
	}{
		"Success": {
			true,
			func(info *FileInfo) error {
				return nil
			},
			nil,
		},
		"Nil ReaderCloser": {
			false,
			func(info *FileInfo) error {
				return nil
			},
			"nil zip.ReadCloser",
		},
		"Walk Error": {
			true,
			func(info *FileInfo) error {
				return fmt.Errorf("error")
			},
			"error",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			r := Repo{
				RepositoryURL: tests.ReleaseRepo,
				ArchiveName:   "verbis_0.0.1_darwin_amd64.zip",
			}

			if test.open {
				err := r.Open()
				assert.NoError(t, err)
				defer r.Close()
			}

			got := r.Walk(test.fn)
			if got != nil {
				assert.Contains(t, got.Error(), test.want)
				return
			}

			assert.Equal(t, test.want, got)
		})
	}
}

func TestRepo_Close(t *testing.T) {
	tt := map[string]struct {
		input Repo
		open bool
		want interface{}
	}{
		"Success": {
			Repo{
				RepositoryURL: tests.ReleaseRepo,
				ArchiveName:   "verbis_0.0.1_darwin_amd64.zip",
			},
			true,
			nil,
		},
		"Remove Error": {
			Repo{tempDir: "."},
			false,
			"error removing temp directory",
		},
		"No Temp Dir": {
			Repo{
				RepositoryURL: "wrong",
				ArchiveName:   "",
			},
			false,
			nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			if test.open {
				err := test.input.Open()
				assert.NoError(t, err)
			}

			got := test.input.Close()
			if got != nil {
				assert.Contains(t, got.Error(), test.want)
				return
			}

			assert.Equal(t, test.want, got)
		})
	}
}

func TestRepo_GetTags(t *testing.T) {
	tt := map[string]struct {
		input Information
		want interface{}
	}{
		"Success": {
			Information{Owner: "cli", Name:  "cli"},
			nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			r := Repo{info: test.input}
			got, err := r.getTags()
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}
			assert.NotEmpty(t, got)
		})
	}
}