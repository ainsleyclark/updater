// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"database/sql"
	"io/ioutil"
)

// runMigrations sorts the migrations and loops over them.
// If there is a migration to run it will be processed
// and committed if the database exists.
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
		shouldRun := u.version.LessThan(migration.toSemVer())
		if !shouldRun {
			continue
		}

		code, err := u.process(migration, tx)
		if err != nil {
			rollBackErr := u.rollBack(tx, down)
			if rollBackErr != nil {
				// In a dirty state
				return code, rollBackErr
			}
			return code, err
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

// rollback reverse the changes from the database (if
// there is one) and the callbacks.
func (u *Updater) rollBack(tx *sql.Tx, down []CallBackFn) error {
	if u.opts.hasDB {
		err := tx.Rollback()
		if err != nil {
			return err
		}
	}

	for _, fn := range down {
		err := fn()
		if err != nil {
			return err
		}
	}

	return nil
}

// process reads the migration and executes the migration
// if there is one. Calls the callback function if there
// is one set.
func (u *Updater) process(m *Migration, tx *sql.Tx) (Status, error) {
	migration, err := ioutil.ReadAll(m.SQL)
	if err != nil {
		return Unknown, err
	}

	if u.opts.hasDB {
		_, err = tx.Exec(string(migration))
		if err != nil {
			return DatabaseError, err
		}
	}

	if m.hasCallBack() {
		err := m.CallBackUp()
		if err != nil {
			return CallBackError, err
		}
	}

	return Updated, nil
}
