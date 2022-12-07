// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// Settings for the plugin.
type Settings struct {
	Source        string
	Destination   string
	Authorization string
	Username      string
	Password      string
	MD5           string
	SHA256        string

	destination string
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	// Verify the source url
	source := p.settings.Source
	if source == "" {
		return fmt.Errorf("no source provided")
	}

	u, err := url.Parse(source)
	if err != nil {
		return fmt.Errorf("could not parse url %s: %w", source, err)
	}

	// Verify the destination
	destination := filepath.ToSlash(p.settings.Destination)
	if destination == "" {
		destination = path.Base(u.Path)
	} else if strings.HasSuffix(destination, "/") {
		destination = path.Join(destination, path.Base(u.Path))
	}

	destination = filepath.FromSlash(path.Clean(destination))
	err = os.MkdirAll(filepath.Dir(destination), os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating directory failed: %w", err)
	}
	p.settings.destination = destination

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	req, err := http.NewRequestWithContext(
		p.network.Context,
		"GET",
		p.settings.Source,
		nil,
	)

	if err != nil {
		return fmt.Errorf("initializing request failed: %w", err)
	}

	p.addAuth(req)

	client := p.network.Client
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		p.addAuth(req)
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"source":      p.settings.Source,
		"destination": p.settings.destination,
	}).Info("Downloading file")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("executing request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed, status %s", http.StatusText(resp.StatusCode))
	}

	target, err := os.Create(p.settings.destination)
	if err != nil {
		return fmt.Errorf("creating destination failed: %w", err)
	}
	defer target.Close()

	_, err = io.Copy(target, resp.Body)
	if err != nil {
		return fmt.Errorf("copying destination failed: %w", err)
	}

	var h hash.Hash
	alg := ""
	exp := ""

	if p.settings.SHA256 != "" {
		exp = p.settings.SHA256
		alg = "SHA256"
		h = sha256.New()
	} else if p.settings.MD5 != "" {
		exp = p.settings.MD5
		alg = "MD5"
		h = md5.New()
	}

	if exp != "" {
		logrus.WithField("hash", alg).Info("Computing checksum")
		_, _ = target.Seek(0, 0)

		if _, err := io.Copy(h, target); err != nil {
			defer os.Remove(target.Name())
			return fmt.Errorf("failed to compare checksum: %w", err)
		}

		check := fmt.Sprintf("%x", h.Sum(nil))

		if exp != check {
			defer os.Remove(target.Name())
			return fmt.Errorf("checksum doesn't match, got %s and expected %s", check, exp)
		}
		logrus.WithField("checksum", check).Info("Checksum matched")
	}

	return nil
}

func (p *Plugin) addAuth(req *http.Request) {
	if p.settings.Username != "" && p.settings.Password != "" {
		req.SetBasicAuth(p.settings.Username, p.settings.Password)
	}
	if p.settings.Authorization != "" {
		req.Header.Add("Authorization", p.settings.Authorization)
	}
}
