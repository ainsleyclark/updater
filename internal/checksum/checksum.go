// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checksum

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrMismatch is returned by compare when no checksum
	// could be found from the URL.
	ErrMismatch = errors.New("checksum didn't match")
	// ErrNoCheckSum is returned by compare when no checksum
	// could be found from the URL.
	ErrNoCheckSum = errors.New("no checksum for archive found")
	// Separator describes the delimiter to separate checksums
	// and archive names from GitHub.
	Separator = "  "
)

// Compare retrieves a checksum files from git and compares
// the hash between the sum on the internet and the sum
// described by the path input.
func Compare(url, path string) (err error) {
	remoteSum, err := getRemoteSum(url, path)
	if err != nil {
		return
	}
	localSum, err := fileSHA256(path)
	if err != nil {
		return
	}
	if remoteSum != localSum {
		return ErrMismatch
	}
	return nil
}

// getRemoteSum retrieves the checksums file from Git and
// separates the release versions (os) by scanning
// each line. The archive name is compared to
// find a match.
func getRemoteSum(url, path string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), Separator)
		if len(line) != 2 { //nolint:gomnd
			continue
		}
		if line[1] == filepath.Base(path) {
			return line[0], nil
		}
	}

	return "", ErrNoCheckSum
}

// fileSHA256 returns the hash of a file which is parsed by
// the given path input.
func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
