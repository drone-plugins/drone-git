package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var version = "1.0.0" // auto-generated version number

func main() {
	app := cli.NewApp()
	app.Name = "git"
	app.Usage = "git clone plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "remote",
			Usage:  "git remote url",
			EnvVar: "CI_REMOTE_URL",
		},
		cli.StringFlag{
			Name:   "path",
			Usage:  "git clone path",
			EnvVar: "CI_WORKSPACE",
		},
		cli.StringFlag{
			Name:   "sha",
			Usage:  "git commit sha",
			EnvVar: "CI_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "CI_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "CI_BUILD_EVENT",
		},
		cli.StringFlag{
			Name:   "number",
			Usage:  "build number",
			EnvVar: "CI_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "netrc.machine",
			Usage:  "netrc machine",
			EnvVar: "CI_NETRC_MACHINE",
		},
		cli.StringFlag{
			Name:   "netrc.username",
			Usage:  "netrc username",
			EnvVar: "CI_NETRC_USERNAME",
		},
		cli.StringFlag{
			Name:   "netrc.password",
			Usage:  "netrc password",
			EnvVar: "CI_NETRC_PASSWORD",
		},
		cli.IntFlag{
			Name:   "depth",
			Usage:  "clone depth",
			EnvVar: "PLUGIN_RECURSIVE",
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
			EnvVar: "PLUGIN_SUBMODULES_UPDATE_REMOTE",
		},
		cli.GenericFlag{
			Name:   "submodule-override",
			Usage:  "json map of submodule overrides",
			EnvVar: "PLUGIN_SUBMODULE_OVERRIDE",
			Value:  &MapFlag{},
		},
	}
	app.Run(os.Args)

}

func run(c *cli.Context) {
	plugin := Plugin{
		Repo: Repo{
			Clone: c.String("remote"),
		},
		Build: Build{
			Number: c.Int("number"),
			Commit: c.String("sha"),
			Event:  c.String("event"),
			Path:   c.String("path"),
			Ref:    c.String("ref"),
		},
		Netrc: Netrc{
			Login:    c.String("netrc.login"),
			Machine:  c.String("netrc.machine"),
			Password: c.String("netrc.password"),
		},
		Config: Config{
			Depth:           c.Int("depth"),
			Recursive:       c.BoolT("recursive"),
			SkipVerify:      c.Bool("skip-verify"),
			SubmoduleRemote: c.Bool("submodule-update-remote"),
			Submodules:      c.Generic("submodule-override").(*MapFlag).Get(),
		},
	}

	if err := plugin.Exec(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
