package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "0.0.0"
	build   = "0"
)

func main() {
	app := cli.NewApp()
	app.Name = "git plugin"
	app.Usage = "git plugin"
	app.Version = fmt.Sprintf("%s+%s", version, build)
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "remote",
			Usage:  "git remote url",
			EnvVar: "PLUGIN_REMOTE,DRONE_REMOTE_URL",
		},
		cli.StringFlag{
			Name:   "path",
			Usage:  "git clone path",
			EnvVar: "PLUGIN_PATH,DRONE_WORKSPACE",
		},
		cli.StringFlag{
			Name:   "sha",
			Usage:  "git commit sha",
			EnvVar: "PLUGIN_SHA,DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "PLUGIN_REF,DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.StringFlag{
			Name:   "netrc.machine",
			Usage:  "netrc machine",
			EnvVar: "DRONE_NETRC_MACHINE",
		},
		cli.StringFlag{
			Name:   "netrc.username",
			Usage:  "netrc username",
			EnvVar: "DRONE_NETRC_USERNAME",
		},
		cli.StringFlag{
			Name:   "netrc.password",
			Usage:  "netrc password",
			EnvVar: "DRONE_NETRC_PASSWORD",
		},
		cli.IntFlag{
			Name:   "depth",
			Usage:  "clone depth",
			EnvVar: "PLUGIN_DEPTH",
		},
		cli.BoolTFlag{
			Name:   "recursive",
			Usage:  "clone submodules",
			EnvVar: "PLUGIN_RECURSIVE",
		},
		cli.BoolFlag{
			Name:   "tags",
			Usage:  "clone tags",
			EnvVar: "PLUGIN_TAGS",
		},
		cli.BoolFlag{
			Name:   "skip-verify",
			Usage:  "skip tls verification",
			EnvVar: "PLUGIN_SKIP_VERIFY",
		},
		cli.BoolFlag{
			Name:   "submodule-update-remote",
			Usage:  "update remote submodules",
			EnvVar: "PLUGIN_SUBMODULES_UPDATE_REMOTE,PLUGIN_SUBMODULE_UPDATE_REMOTE",
		},
		cli.GenericFlag{
			Name:   "submodule-override",
			Usage:  "json map of submodule overrides",
			EnvVar: "PLUGIN_SUBMODULE_OVERRIDE",
			Value:  &MapFlag{},
		},
		cli.DurationFlag{
			Name:   "backoff",
			Usage:  "backoff duration",
			EnvVar: "PLUGIN_BACKOFF",
			Value:  5 * time.Second,
		},
		cli.IntFlag{
			Name:   "backoff-attempts",
			Usage:  "backoff attempts",
			EnvVar: "PLUGIN_ATTEMPTS",
			Value:  5,
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}

}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Repo: Repo{
			Clone: c.String("remote"),
		},
		Build: Build{
			Commit: c.String("sha"),
			Event:  c.String("event"),
			Path:   c.String("path"),
			Ref:    c.String("ref"),
		},
		Netrc: Netrc{
			Login:    c.String("netrc.username"),
			Machine:  c.String("netrc.machine"),
			Password: c.String("netrc.password"),
		},
		Config: Config{
			Depth:           c.Int("depth"),
			Tags:            c.Bool("tags"),
			Recursive:       c.BoolT("recursive"),
			SkipVerify:      c.Bool("skip-verify"),
			SubmoduleRemote: c.Bool("submodule-update-remote"),
			Submodules:      c.Generic("submodule-override").(*MapFlag).Get(),
		},
		Backoff: Backoff{
			Attempts: c.Int("backoff-attempts"),
			Duration: c.Duration("backoff"),
		},
	}

	return plugin.Exec()
}
