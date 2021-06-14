// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

func (g *Repo) checkSum(tag string) ([]byte, error) {
	const op = "Repo.CheckSum"

	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/checksums.txt", "ainsleyclark/verbis", tag)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "  ")
		if len(line) != 2 {
			continue
		}
		if line[1] == g.ArchiveName {
			return []byte(line[0]), nil
		}
	}

	return nil, ErrNoCheckSum
}
