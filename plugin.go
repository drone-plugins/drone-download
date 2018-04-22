package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/jackspirou/syscerts"
	"github.com/pkg/errors"
)

type (
	Config struct {
		Source      string
		Destination string
		Username    string
		Password    string
		SkipVerify  bool
		MD5         string
		SHA265      string
	}

	Plugin struct {
		Config Config
	}
)

func (p Plugin) Exec() error {
	destination := p.Config.Destination

	if destination == "" {
		u, err := url.Parse(p.Config.Source)

		if err != nil {
			return errors.Wrap(err, "parsing source failed")
		}

		destination = path.Base(u.Path)
	}

	log.Printf("downloading to %s", destination)

	client := &http.Client{
		Timeout: time.Minute * 5,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				RootCAs:            syscerts.SystemRootsPool(),
				InsecureSkipVerify: p.Config.SkipVerify,
			},
		},
	}

	req, err := http.NewRequest(
		"GET",
		p.Config.Source,
		nil,
	)

	if err != nil {
		return errors.Wrap(err, "initializing request failed")
	}

	if p.Config.Username != "" && p.Config.Password != "" {
		req.SetBasicAuth(p.Config.Username, p.Config.Password)
	}

	resp, err := client.Do(req)

	if err != nil {
		return errors.Wrap(err, "executing request failed")
	}

	defer resp.Body.Close()

	target, err := os.Create(destination)

	if err != nil {
		return errors.Wrap(err, "creating destination failed")
	}

	defer target.Close()

	if p.Config.MD5 != "" {
		h := md5.New()

		if _, err := io.Copy(h, resp.Body); err != nil {
			return errors.Wrap(err, "failed to compare checksum")
		}

		check := fmt.Sprintf("%x", h.Sum(nil))

		if p.Config.MD5 != check {
			return fmt.Errorf("checksum doesn't match, got %s and expected %s", check, p.Config.MD5)
		}
	}

	if p.Config.SHA265 != "" {
		h := sha256.New()

		if _, err := io.Copy(h, resp.Body); err != nil {
			return errors.Wrap(err, "failed to compare checksum")
		}

		check := fmt.Sprintf("%x", h.Sum(nil))

		if p.Config.MD5 != check {
			return fmt.Errorf("checksum doesn't match, got %s and expected %s", check, p.Config.SHA265)
		}
	}

	_, err = io.Copy(target, resp.Body)

	if err != nil {
		return errors.Wrap(err, "copying destination failed")
	}

	return nil
}
