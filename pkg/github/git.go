// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

// |||||||||||||||||||||||||||||||||||||||||||||||||||||||

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ainsleyclark/updater/pkg/checksum"
	"github.com/ainsleyclark/updater/pkg/fileio"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type Provider interface {
	Open() (err error)
	Walk(walkFn WalkFunc) error
	Close() error
}

type Repo struct {
	RepositoryURL string
	ArchiveName   string
	ChecksumName  string
	tempDir       string
	archivePath   string
	info          Information
	reader        *zip.ReadCloser
}

// tag defines the data used to unmarshal the response from
// github. `Name` is only required to compare versions
// and obtain archive information.
type tag struct {
	Name string `json:"name"`
}

// Information contains the name and owner of the repository
// used for obtaining tags, archive URL and checksums
// if they are attached.
type Information struct {
	Owner string
	Name  string
}

// getInfo returns the owner and name of the github url
// using regex. An error will be returned if the url
// is invalid.
func (r *Repo) getInfo() error {
	re := regexp.MustCompile(`github\.com/(.*?)/(.*?)$`)
	submatches := re.FindAllStringSubmatch(r.RepositoryURL, 1)
	if len(submatches) < 1 {
		return errors.New("invalid github URL:" + r.RepositoryURL)
	}
	r.info = Information{
		Owner: submatches[0][1],
		Name:  submatches[0][2],
	}
	return nil
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

func (r *Repo) Download() (string, error) {
	version, err := r.LatestVersion()
	if err != nil {
		return "", err
	}

	// Retrieve the zip folder
	resp, err := http.Get(r.getDownloadUrl(version, r.ArchiveName))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	r.tempDir, err = fileio.TempDirectory()
	if err != nil {
		return "", err
	}

	r.archivePath = filepath.Join(r.tempDir, r.ArchiveName)
	archiveFile, err := os.Create(r.archivePath)
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
		err := checksum.Compare(r.getDownloadUrl(version, r.ChecksumName), r.archivePath)
		if err != nil {
			return "", err
		}
	}

	return r.archivePath, nil
}



func (r *Repo) Open() error {
	version, err := r.LatestVersion()
	if err != nil {
		return err
	}

	// Retrieve the zip folder
	resp, err := http.Get(r.getDownloadUrl(version, r.ArchiveName))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	r.tempDir, err = fileio.TempDirectory()
	if err != nil {
		return err
	}

	r.archivePath = filepath.Join(r.tempDir, r.ArchiveName)
	archiveFile, err := os.Create(r.archivePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(archiveFile, resp.Body)
	archiveFile.Close()
	if err != nil {
		return err
	}

	// Compare checksums if a name is set.
	if r.ChecksumName != "" {
		err := checksum.Compare(r.getDownloadUrl(version, r.ChecksumName), r.archivePath)
		if err != nil {
			return err
		}
	}

	r.reader, err = zip.OpenReader(r.archivePath)
	if err != nil {
		return err
	}

	return nil
}

// WalkFunc is used for walking over the repository and
// collecting file info by iterating over the zip
// file.
type WalkFunc func(info *FileInfo) error

// FileInfo defines the information sent back from the walk
// function.
type FileInfo struct {
	Path     string
	Mode     os.FileMode
	Modified time.Time
}

// Walk iterates over the zip folder stored in a temporary
// directory. If the zip does not exist or the zip
// ReadCloser is nil, and error will be returned.
func (r *Repo) Walk(walkFn WalkFunc) error {
	if r.reader == nil {
		return errors.New("nil zip.ReadCloser")
	}

	for _, f := range r.reader.File {
		if f != nil {
			err := walkFn(&FileInfo{
				Path:     f.Name,
				Mode:     f.Mode(),
				Modified: f.Modified,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Repo) findZipFile(path string) (*zip.File, error) {
	for _, f := range r.reader.File {
		if path == f.Name {
			return f, nil
		}
	}
	return nil, fmt.Errorf("no zip file found with the path: %s", path)
}

func (r *Repo) Copy(src, dest string) error {
	zipFile, err := r.findZipFile(src)
	if err != nil {
		return err
	}

	file, err := zipFile.Open()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}

	return nil
}

// Close removes the temporary directory used to store the
// zip folder downloaded from Github. It then closes the
// zip.ReaderCloser attached.
func (r *Repo) Close() error {
	if r.tempDir == "" {
		return nil
	}

	err := os.RemoveAll(r.tempDir)
	if err != nil {
		return errors.New("error removing temp directory: " + err.Error())
	}
	r.tempDir = ""

	return r.reader.Close()
}

// getDownloadUrl returns the URL of a download from the
// repository based on the input tag name, and the name
// of the archive (could be a zip or checksums.txt).
func (r *Repo) getDownloadUrl(tag string, name string) string {
	return fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", r.info.Owner, r.info.Name, tag, name)
}

// getTags retrieves the latest tag information from
// GitHub and returns a slice of tags containing
// the name of the release.
func (r *Repo) getTags() ([]tag, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", r.info.Owner, r.info.Name))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tags []tag
	err = json.NewDecoder(resp.Body).Decode(&tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
