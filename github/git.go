// Copyright 2020 The Verbis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type Provider interface {
	Open() (err error)
	Walk(walkFn WalkFunc) error
	Close() error
}

type Repo struct {
	RepoURL     string
	ArchiveName string
	tempDir     string
	archivePath string
	reader      *zip.ReadCloser
}

// tag struct used to unmarshal response from Repo
// https://api.github.com/repos/ownerName/projectName/tags
type tag struct {
	Name string `json:"name"`
}

// A FileInfo describes a file given by a provider
type FileInfo struct {
	Path string
	Mode os.FileMode
}

// WalkFunc is the type of the function called for each file or directory
// visited by Walk.
// path is relative
type WalkFunc func(info *FileInfo) error

var (
	ErrNoCheckSum = errors.New("no checksum for archive found")
	tagsUrl       = "https://api.Repo.com/repos/ainsleyclark/verbis/tags"
)

const (
	TempDirName = "verbis-updater"
)

func (g *Repo) Open() (err error) {
	version, err := g.LatestVersion()
	if err != nil {
		return
	}

	archiveURL := g.getArchiveURL(version)
	resp, err := http.Get(archiveURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	g.tempDir, err = ioutil.TempDir("", TempDirName)
	if err != nil {
		return
	}

	g.archivePath = filepath.Join(g.tempDir, g.ArchiveName)
	archiveFile, err := os.Create(g.archivePath)
	if err != nil {
		return
	}

	_, err = io.Copy(archiveFile, resp.Body)
	archiveFile.Close()
	if err != nil {
		return
	}

	g.reader, err = zip.OpenReader(g.archivePath)
	if err != nil {
		return
	}

	return
}

// Walk
func (g *Repo) Walk(walkFn WalkFunc) error {
	if g.reader == nil {
		return errors.New("nil zip.ReadCloser")
	}

	for _, f := range g.reader.File {
		if f != nil {
			err := walkFn(&FileInfo{
				Path: f.Name,
				Mode: f.Mode(),
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Repo) Close() error {
	if len(g.tempDir) <= 0 {
		return nil
	}

	err := os.RemoveAll(g.tempDir)
	if err != nil {
		return err
	}
	g.tempDir = ""

	return g.reader.Close()
}

// GetArchiveURL
func (g *Repo) getArchiveURL(tag string) string {
	return fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", "ainleyclark/verbis", tag, g.ArchiveName)
}

// getTags
func (g *Repo) getTags() ([]tag, error) {
	//tagsUrl := "https://api.Github.com/repos/" + api.Repo + "/verbis/tags"

	resp, err := http.Get(tagsUrl)
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

// LatestVersion
func (g *Repo) LatestVersion() (string, error) {
	tags, err := g.getTags()
	if err != nil {
		return "", err
	}

	if len(tags) < 1 {
		return "", errors.New("repo has no tags")
	}

	return tags[0].Name, nil
}