// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater

const (
	v001 = "UPDATE my_table SET name = 'tom' WHERE id = 1"
)

//
//func TestUpdater_Run(t *testing.T) {
//	tt := map[string]struct {
//		input migrationRegistry
//		mock  func(m sqlmock.Sqlmock)
//		want  interface{}
//		code int
//	}{
//		"Simple": {
//			migrationRegistry{
//				&Migration{Version: "v0.0.1", MigrationPath: "v0.0.1.sql", Stage: Major},
//			},
//			func(m sqlmock.Sqlmock) {
//				m.ExpectBegin()
//				m.ExpectExec(v001).
//					WillReturnResult(sqlmock.NewResult(1, 1))
//				m.ExpectCommit()
//			},
//			nil,
//			4,
//		},
//		"Wrong File Name": {
//			migrationRegistry{
//				&Migration{Version: "v0.0.0", MigrationPath: "v0.0.1-wrong.sql", Stage: Major},
//			},
//			func(m sqlmock.Sqlmock) {
//				m.ExpectBegin()
//				m.ExpectRollback()
//			},
//			"open v0.0.1-wrong.sql: file does not exist",
//		},
//		"RollBack Error": {
//			migrationRegistry{
//				&Migration{Version: "v0.0.0", MigrationPath: "v0.0.1-wrong.sql", Stage: Major},
//			},
//			func(m sqlmock.Sqlmock) {
//				m.ExpectBegin()
//				m.ExpectRollback().
//					WillReturnError(fmt.Errorf("error"))
//			},
//			"error",
//		},
//		//"Begin Error": {
//		//	nil,
//		//	func(m sqlmock.Sqlmock) {
//		//		m.ExpectBegin().
//		//			WillReturnError(fmt.Errorf("error"))
//		//	},
//		//	"error",
//		//},
//	}
//
//	for name, test := range tt {
//		t.Run(name, func(t *testing.T) {
//			db, mock, err := sqlmock.New()
//			if err != nil {
//				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//			}
//
//			defer func() {
//				migrations = make(migrationRegistry, 0)
//				defer db.Close()
//			}()
//
//			migrations = test.input
//			u, err := New(Options{
//				DB:            db,
//				Embed:         testdata.Static,
//				Version:       "0.0.0",
//				RepositoryURL: "https://github.com/ainsleyclark/verbis",
//			})
//			assert.NoError(t, err)
//
//			test.mock(mock)
//
//			err = u.Run()
//			if err != nil {
//				color.Red.Println(err)
//				assert.Contains(t, err.Error(), test.want)
//				return
//			}
//
//			err = mock.ExpectationsWereMet()
//			if err != nil {
//				t.Errorf("there were unfulfilled expectations: %s", err)
//			}
//		})
//	}
//}
