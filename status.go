// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import "github.com/mouuff/go-rocket-update/pkg/updater"

// Status defines the status codes returned by the Update()
// function used for debugging any issues with updating
// the application.
type Status int

const (
	// Unknown update status (something went wrong)
	Unknown Status = iota
	// DatabaseError is returned by update when a database
	// connection could not be established or there was
	// an error processing the transaction.
	DatabaseError = 1
	// ExecutableError is returned by update when there was
	// a error updating the executable from GitHub.
	ExecutableError = 2
	// CallBackError is returned by update when there was
	// a error with one of the migration callbacks.
	CallBackError = 3
	// UpToDate status is used to define when the application
	// is already up to date.
	UpToDate = 5
	// Updated is the success status code returned by Update
	// when everything passed.
	Updated = 6
)

// getExecStatus transforms the pkg updater status into
// the Status codes listed above.
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
