// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

// DO NOT MODIFY THIS FILE DIRECTLY

package main

import (
	"os"

	"github.com/drone-plugins/drone-download/plugin"
	"github.com/drone-plugins/drone-plugin-lib/errors"
	"github.com/drone-plugins/drone-plugin-lib/urfave"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var version = "unknown"

func main() {
	settings := &plugin.Settings{}

	if _, err := os.Stat("/run/drone/env"); err == nil {
		_ = godotenv.Overload("/run/drone/env")
	}

	app := &cli.App{
		Name:    "drone-download",
		Usage:   "download a file",
		Version: version,
		Flags:   append(settingsFlags(settings), urfave.Flags()...),
		Action:  run(settings),
	}

	if err := app.Run(os.Args); err != nil {
		errors.HandleExit(err)
	}
}

func run(settings *plugin.Settings) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		urfave.LoggingFromContext(ctx)

		plugin := plugin.New(
			*settings,
			urfave.PipelineFromContext(ctx),
			urfave.NetworkFromContext(ctx),
		)

		if err := plugin.Validate(); err != nil {
			if e, ok := err.(errors.ExitCoder); ok {
				return e
			}

			return errors.ExitMessagef("validation failed: %w", err)
		}

		if err := plugin.Execute(); err != nil {
			if e, ok := err.(errors.ExitCoder); ok {
				return e
			}

			return errors.ExitMessagef("execution failed: %w", err)
		}

		return nil
	}
}

func settingsFlags(settings *plugin.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "source",
			Usage:       "source url for the download",
			EnvVars:     []string{"PLUGIN_SOURCE"},
			Destination: &settings.Source,
		},
		&cli.StringFlag{
			Name:        "destination",
			Usage:       "destination for the download",
			EnvVars:     []string{"PLUGIN_DESTINATION"},
			Destination: &settings.Destination,
		},
		&cli.StringFlag{
			Name:        "authorization",
			Usage:       "value to send in the authorization header",
			EnvVars:     []string{"PLUGIN_AUTHORIZATION", "DOWNLOAD_AUTHORIZATION"},
			Destination: &settings.Authorization,
		},
		&cli.StringFlag{
			Name:        "username",
			Usage:       "username for basic auth",
			EnvVars:     []string{"PLUGIN_USERNAME", "DOWNLOAD_USERNAME"},
			Destination: &settings.Username,
		},
		&cli.StringFlag{
			Name:        "password",
			Usage:       "password for basic auth",
			EnvVars:     []string{"PLUGIN_PASSWORD", "DOWNLOAD_PASSWORD"},
			Destination: &settings.Password,
		},
		&cli.StringFlag{
			Name:        "md5-checksum",
			Usage:       "checksum in md5 format",
			EnvVars:     []string{"PLUGIN_MD5"},
			Destination: &settings.MD5,
		},
		&cli.StringFlag{
			Name:        "sha256-checksum",
			Usage:       "checksum in sha256 format",
			EnvVars:     []string{"PLUGIN_SHA256", "PLUGIN_SHA265"},
			Destination: &settings.SHA256,
		},
	}
}
