// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/ainsleyclark/updater/pkg/github"
	"github.com/ainsleyclark/updater/pkg/updater"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	exec, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	base := filepath.Dir(exec)

	u := updater.Updater{
		Github: github.Repo{
			RepositoryURL: "https://github.com/ainsleyclark/verbis",
			ArchiveName:   fmt.Sprintf("verbis_%s_%s_%s.zip", "0.0.1", runtime.GOOS, runtime.GOARCH),
			ChecksumName:  "checksums.txt",
		},
		Files: updater.Files{
			{RemotePath: "verbis/build/admin/**/css", LocalPath: filepath.Join(base, "admin")},
			{RemotePath: "verbis/build/admin/**/**", LocalPath: filepath.Join(base, "admin")},
			{RemotePath: "verbis/build/admin/index.html", LocalPath: filepath.Join(base, "index.html")},
		},
		Version:              "",
		BackupExtension:      "",
		RemoteExecutablePath: "verbis/build/verbis",
	}



	fmt.Println(u.Update())
}
