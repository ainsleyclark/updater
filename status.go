// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import "github.com/mouuff/go-rocket-update/pkg/updater"

type Status int

const (
	// Unknown update status (something went wrong)
	Unknown         Status = iota
	DatabaseError          = 1
	ExecutableError        = 2
	UpToDate               = 3
	Updated                = 4
)

func getExecStatus(status updater.UpdateStatus) Status {
	switch status {
	case updater.Unknown:
		return ExecutableError
	case updater.UpToDate:
		return UpToDate
	case updater.Updated:
		return Updated
	}
	return Unknown
}
