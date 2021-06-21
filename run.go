// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

import (
	"database/sql"
	"github.com/gookit/color"
)

func (u *Updater) runMigrations() (Status, error) {
	tx, err := u.opts.DB.Begin()
	if err != nil {
		return DatabaseError, err
	}

	// TODO, what happens if there is no database but a callback?

	migrations.Sort()

	var down []CallBackFn
	for _, migration := range migrations {
		shouldRun := u.version.LessThanOrEqual(migration.toSemVer())
		if !shouldRun {
			continue
		}

		err := u.process(migration, tx)
		color.Blue.Println(err)
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

	err = tx.Commit()
	if err != nil {
		return DatabaseError, err
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
	migration, err := u.opts.Embed.ReadFile(m.MigrationPath)
	if err != nil {
		return err
	}

	_, err = tx.Exec(string(migration))
	if err != nil {
		return err
	}

	if m.hasCallBack() {
		err := m.CallBackUp()
		if err != nil {
			return err
		}
	}

	return nil
}
