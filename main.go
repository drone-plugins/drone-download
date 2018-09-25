package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var (
	version = "0.0.0"
	build   = "0"
)

func main() {
	app := cli.NewApp()
	app.Name = "download plugin"
	app.Usage = "download plugin"
	app.Action = run
	app.Version = fmt.Sprintf("%s+%s", version, build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "source",
			Usage:  "source url for the download",
			EnvVar: "PLUGIN_SOURCE",
		},
		cli.StringFlag{
			Name:   "destination",
			Usage:  "destination for the download",
			EnvVar: "PLUGIN_DESTINATION",
		},
		cli.StringFlag{
			Name:   "authorization",
			Usage:  "value to send in the authorization header",
			EnvVar: "PLUGIN_AUTHORIZATION,DOWNLOAD_AUTHORIZATION",
		},
		cli.StringFlag{
			Name:   "username",
			Usage:  "username for basic auth",
			EnvVar: "PLUGIN_USERNAME,DOWNLOAD_USERNAME",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "password for basic auth",
			EnvVar: "PLUGIN_PASSWORD,DOWNLOAD_PASSWORD",
		},
		cli.BoolFlag{
			Name:   "skip-verify",
			Usage:  "skip ssl verification",
			EnvVar: "PLUGIN_SKIP_VERIFY",
		},
		cli.StringFlag{
			Name:   "md5-checksum",
			Usage:  "checksum in md5 format",
			EnvVar: "PLUGIN_MD5",
		},
		cli.StringFlag{
			Name:   "sha265-checksum",
			Usage:  "checksum in sha265 format",
			EnvVar: "PLUGIN_SHA265",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Config: Config{
			Source:        c.String("source"),
			Destination:   c.String("destination"),
			Authorization: c.String("authorization"),
			Username:      c.String("username"),
			Password:      c.String("password"),
			SkipVerify:    c.Bool("skip-verify"),
			MD5:           c.String("md5-checksum"),
			SHA265:        c.String("sha265-checksum"),
		},
	}

	if plugin.Config.Source == "" {
		return errors.New("Missing source URL")
	}

	return plugin.Exec()
}
