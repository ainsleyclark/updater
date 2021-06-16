// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package patcher

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestPatcher_Apply(t *testing.T) {
	tt := map[string]struct {
		fn   func() *Patcher
		want interface{}
	}{
		"Empty": {
			func() *Patcher {
				return &Patcher{}
			},
			nil,
		},
		"Read File Error": {
			func() *Patcher {
				return &Patcher{
					files: []*File{{SourcePath: "", DestinationPath: "", Mode: 0}},
				}
			},
			"no such file or directory",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			p := test.fn()
			err := p.apply()
			if err != nil {
				require.Contains(t, err.Error(), test.want)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestPatcher_AddFile(t *testing.T) {
	tt := map[string]struct {
		patcher *Patcher
		input   *File
		want    interface{}
	}{
		"Empty": {
			&Patcher{},
			&File{SourcePath: "test"},
			1,
		},
		"Nil": {
			&Patcher{},
			nil,
			0,
		},
		"With File": {
			&Patcher{
				files: []*File{
					{SourcePath: "test"},
				},
			},
			&File{SourcePath: "test"},
			2,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			test.patcher.AddFile(test.input)
			require.Equal(t, test.want, len(test.patcher.files))
		})
	}
}

func TestPatcher_Backup(t *testing.T) {
	tmp, err := ioutil.TempFile("", "updater-test")
	require.NoError(t, err)
	backup := tmp.Name() + ".test"

	defer func() {
		os.Remove(tmp.Name())
		os.RemoveAll(backup)
	}()

	p := &Patcher{
		DestinationPath: tmp.Name(),
		BackupPath:      backup,
		backupPossible:  true,
	}

	got := p.backup()
	require.NoError(t, got)

	_, err = os.Stat(backup)
	require.NoError(t, got)
}

func TestPatcher_Backup_Nil(t *testing.T) {
	p := &Patcher{
		backupPossible: false,
	}
	got := p.backup()
	require.NoError(t, got)
}
