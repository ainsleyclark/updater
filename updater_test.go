// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"fmt"
	"github.com/mouuff/go-rocket-update/pkg/provider"
	"github.com/mouuff/go-rocket-update/pkg/updater"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	tt := map[string]struct {
		input Options
		error bool
	}{
		"Success": {
			Options{GithubURL: "https://github.com/ainsleyclark/verbis", Version: "0.0.1"},
			false,
		},
		"Bad Version": {
			Options{GithubURL: "https://github.com/ainsleyclark/verbis", Version: "wrong"},
			true,
		},
		"Bad Options": {
			Options{},
			true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := New(test.input)
			if test.error && err != nil {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, test.input, got.opts)
			assert.Equal(t, test.input.Version, got.pkg.Version)
		})
	}
}

var (
	TestErr     = fmt.Errorf("error") //nolint
	TestVersion = "v0.0.1"            //nolint
)

type mockAccessProvider struct{}

func (a *mockAccessProvider) Walk(walkFn provider.WalkFunc) error {
	return nil
}

func (a *mockAccessProvider) Open() error {
	return nil
}

func (a *mockAccessProvider) Close() error {
	return nil
}

func (a *mockAccessProvider) GetLatestVersion() (string, error) {
	return TestVersion, nil
}

func (a *mockAccessProvider) Retrieve(srcPath, destPath string) error {
	return nil
}

type mockAccessProviderErr struct{}

func (a *mockAccessProviderErr) Walk(walkFn provider.WalkFunc) error {
	return nil
}

func (a *mockAccessProviderErr) Open() error {
	return nil
}

func (a *mockAccessProviderErr) Close() error {
	return nil
}

func (a *mockAccessProviderErr) GetLatestVersion() (string, error) {
	return "", fmt.Errorf("error")
}

func (a *mockAccessProviderErr) Retrieve(srcPath, destPath string) error {
	return nil
}

func TestUpdater_HasUpdate(t *testing.T) {
	tt := map[string]struct {
		input provider.Provider
		want  interface{}
	}{
		"Success": {
			&mockAccessProvider{},
			true,
		},
		"Error": {
			&mockAccessProviderErr{},
			TestErr.Error(),
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			u := Updater{pkg: &updater.Updater{Provider: test.input}}
			got, err := u.HasUpdate()
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}
			assert.Equal(t, test.want, got)
		})
	}
}

func TestUpdater_LatestVersion(t *testing.T) {
	tt := map[string]struct {
		input provider.Provider
		want  interface{}
	}{
		"Success": {
			&mockAccessProvider{},
			TestVersion,
		},
		"Error": {
			&mockAccessProviderErr{},
			TestErr.Error(),
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			u := Updater{pkg: &updater.Updater{Provider: test.input}}
			got, err := u.LatestVersion()
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}
			assert.Equal(t, test.want, got)
		})
	}
}

func TestUpdater_Update(t *testing.T) {
	// TODO
}
