// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ainsleyclark/updater/internal/checksum"
	"github.com/ainsleyclark/updater/internal/fileio"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

// Repo defines the Github repository to obtain the update.
// The repository URL must match a github.com account
// and the archive name must be valid. Checksums
// are entirely optional and can be omitted
// if they are not to be validated once
// downloaded.
type Repo struct {
	RepositoryURL string
	ArchiveName   string
	ChecksumName  string
	tempDir       string
	info          information
}

var (
	baseUrl = "https://github.com"
	apiUrl  = "https://api.github.com"
)

// tag defines the data used to unmarshal the response from
// github. `Name` is only required to compare versions
// and obtain archive information.
type tag struct {
	Name string `json:"name"`
}

// information contains the name and owner of the repository
// used for obtaining tags, archive URL and checksums
// if they are attached.
type information struct {
	Owner string
	Name  string
}

// LatestVersion retrieves the latest release (tags) from
// GitHub from the first tag. An error will be returned
// if the repo has no tags or there was a problem
// calling the Git API.
func (r *Repo) LatestVersion() (string, error) {
	err := r.getInfo()
	if err != nil {
		return "", err
	}

	tags, err := r.getTags()
	if err != nil {
		return "", err
	}

	if len(tags) < 1 {
		return "", errors.New("repo has no tags")
	}

	return tags[0].Name, nil
}

// Download retrieves the relevant Github information, latest
// tag and downloads the zip folder to a temporary
// directory. A zip reader is the opened ready
// for walking and copying files and folders.
func (r *Repo) Download() (string, error) {
	version, err := r.LatestVersion()
	if err != nil {
		return "", err
	}

	// Retrieve the zip folder
	resp, err := http.Get(r.getDownloadURL(version, r.ArchiveName))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	r.tempDir, err = fileio.TempDirectory()
	if err != nil {
		return "", err
	}

	zipPath := filepath.Join(r.tempDir, r.ArchiveName)
	archiveFile, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(archiveFile, resp.Body)
	archiveFile.Close()
	if err != nil {
		return "", err
	}

	// Compare checksums if a name is set.
	if r.ChecksumName != "" {
		err := checksum.Compare(r.getDownloadURL(version, r.ChecksumName), zipPath)
		if err != nil {
			return "", err
		}
	}

	return zipPath, nil
}

// Close removes the temporary directory the zip folder
// is stored in. An error will be returned if the
// temporary dir could not be removed.
func (r *Repo) Close() error {
	if r.tempDir == "" {
		return nil
	}
	err := os.RemoveAll(r.tempDir)
	if err != nil {
		return errors.New("error removing temp directory: " + err.Error())
	}
	r.tempDir = ""
	return nil
}

// getInfo returns the owner and name of the github url
// using regex. An error will be returned if the url
// is invalid.
func (r *Repo) getInfo() error {
	re := regexp.MustCompile(`github\.com/(.*?)/(.*?)$`)
	submatches := re.FindAllStringSubmatch(r.RepositoryURL, 1)
	if len(submatches) < 1 {
		return errors.New("invalid github url:" + r.RepositoryURL)
	}
	r.info = information{
		Owner: submatches[0][1],
		Name:  submatches[0][2],
	}
	return nil
}

// getDownloadURL returns the URL of a download from the
// repository based on the input tag name, and the name
// of the archive (could be a zip or checksums.txt).
func (r *Repo) getDownloadURL(tag, name string) string {
	return fmt.Sprintf("%s/%s/%s/releases/download/%s/%s", baseUrl, r.info.Owner, r.info.Name, tag, name)
}

// getTags retrieves the latest tag information from
// GitHub and returns a slice of tags containing
// the name of the release.
func (r *Repo) getTags() ([]tag, error) {
	resp, err := http.Get(fmt.Sprintf("%s/repos/%s/%s/tags", apiUrl, r.info.Owner, r.info.Name))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var tags []tag
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
