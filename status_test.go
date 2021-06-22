// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"github.com/mouuff/go-rocket-update/pkg/updater"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetExecStatus(t *testing.T) {
	tt := map[string]struct {
		input   updater.UpdateStatus
		want    Status
	}{
		"Unknown": {
			updater.Unknown,
			ExecutableError,
		},
		"Up To Date": {
			updater.UpToDate,
			UpToDate,
		},
		"Updated": {
			updater.Updated,
			Updated,
		},
		"Default": {
			updater.UpdateStatus(999),
			Unknown,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := getExecStatus(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
