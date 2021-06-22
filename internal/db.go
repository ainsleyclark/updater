// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"github.com/ainsleyclark/updater"
	"github.com/jmoiron/sqlx"
)

type Migrator interface {
	Migrate(m updater.Migration)
}


type tssr struct {
	tx *sqlx.Tx
}

func (t *tssr) Migrate(m updater.Migration) {

}