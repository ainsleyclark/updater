// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"database/sql"
	"io/ioutil"
)

// db
// migrations
// version
//

func (u *Updater) runMigrations() (Status, error) {
	var (
		err error
		tx *sql.Tx
	)

	if u.opts.hasDB {
		tx, err = u.opts.DB.Begin()
		if err != nil {
			return DatabaseError, err
		}
	}

	migrations.Sort()

	var down []CallBackFn
	for _, migration := range migrations {
		shouldRun := u.version.LessThanOrEqual(migration.toSemVer())
		if !shouldRun {
			continue
		}

		err := u.process(migration, tx)
		if err != nil {
			rollBackErr := u.rollBack(tx, down)
			if rollBackErr != nil {
				// In a dirty state
				return Unknown, rollBackErr
			}
			return Unknown, err
		}

		down = append(down, migration.CallBackDown)
	}

	if u.opts.hasDB {
		err := tx.Commit()
		if err != nil {
			return DatabaseError, err
		}
	}

	return Updated, nil
}

func (u *Updater) rollBack(tx *sql.Tx, down []CallBackFn) error {
	err := tx.Rollback()
	if err != nil {
		return err
	}

	for _, fn := range down {
		err := fn()
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *Updater) process(m *Migration, tx *sql.Tx) error {
	migration, err := ioutil.ReadAll(m.Migration)
	if err != nil {
		return err
	}

	if u.opts.hasDB {
		_, err = tx.Exec(string(migration))
		if err != nil {
			return err
		}
	}

	if m.hasCallBack() {
		err := m.CallBackUp()
		if err != nil {
			return err
		}
	}

	return nil
}
