// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package updater_test

//func TestUpdater(t *testing.T) {
//	u, err := updater.New(updater.Options{
//		RepositoryURL: "https://github.com/ainsleyclark/verbis",
//		Version:       "v0.0.1",
//		Verify:        false,
//		DB:           nil,
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	status, err := u.Update(fmt.Sprintf("verbis_v0.0.2_%s_%s.zip", runtime.GOOS, runtime.GOARCH))
//	if err != nil {
//		return
//	}
//
//	fmt.Println(status)
//}
//
//func init() {
//	err := updater.AddMigration(&updater.SQL{
//		Version:      "v0.0.2",
//		SQL:    strings.NewReader("UPDATE my_table SET 'title' WHERE id = 1"),
//		CallBackUp:   func() error { return nil }, // Runs on up of migration.
//		CallBackDown: func() error { return nil }, // Runs on error of migration.
//		Stage:        updater.Patch,
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//}