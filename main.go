package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/vrischmann/envconfig"
)

const netrcFile = `
machine %s
login %s
password %s
`

type Plugin struct {
	Repo struct {
		Clone string `envconfig:"CI_CLONE_URL"`
	}

	Build struct {
		Path   string `envconfig:"CI_WORKSPACE"`
		Event  string `envconfig:"CI_BUILD_EVENT"`
		Number int    `envconfig:"CI_BUILD_NUMBER"`
		Commit string `envconfig:"CI_COMMIT_SHA"`
		Ref    string `envconfig:"CI_COMMIT_REF"`
	}

	Netrc struct {
		Machine  string `envconfig:"CI_NETRC_MACHINE"`
		Login    string `envconfig:"CI_NETRC_LOGIN"`
		Password string `envconfig:"CI_NETRC_PASSWORD"`
	}

	Config struct {
		Depth           int               `envconfig:"PLUGIN_DEPTH"`
		Recursive       bool              `envconfig:"PLUGIN_RECURSIVE"`
		SkipVerify      bool              `envconfig:"PLUGIN_SKIP_VERIFY"`
		Tags            bool              `envconfig:"PLUGIN_TAGS"`
		Submodules      map[string]string `envconfig:"PLUGIN_SUBMODULE_OVERRIDE"`
		SubmoduleRemote bool              `envconfig:"PLUGIN_SUBMODULE_UPDATE_REMOTE"`
	}
}

var opts = envconfig.Options{
	AllOptional: true,
}

func main() {
	plugin := &Plugin{}

	if err := envconfig.InitWithOptions(&plugin, opts); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := plugin.Exec(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (p Plugin) Exec() error {
	err := os.MkdirAll(p.Build.Path, 0777)
	if err != nil {
		return err
	}

	// generate the .netrc file
	err = writeNetrc(p.Netrc.Login, p.Netrc.Login, p.Netrc.Password)
	if err != nil {
		return err
	}

	var cmds []*exec.Cmd

	if p.Config.SkipVerify {
		cmds = append(cmds, skipVerify())
	}

	// check for a .git directory and whether it's empty
	if isDirEmpty(filepath.Join(p.Build.Path, ".git")) {
		cmds = append(cmds, initGit())
		cmds = append(cmds, remote(p.Repo.Clone))
	}

	switch {
	case isPullRequest(p.Build.Event) || isTag(p.Build.Event, p.Build.Ref):
		cmds = append(cmds, fetch(p.Build.Ref, p.Config.Tags, p.Config.Depth))
		cmds = append(cmds, checkoutHead())
	default:
		cmds = append(cmds, fetch(p.Build.Ref, p.Config.Tags, p.Config.Depth))
		cmds = append(cmds, checkoutSha(p.Build.Commit))
	}

	for name, url := range p.Config.Submodules {
		cmds = append(cmds, remapSubmodule(name, url))
	}

	if p.Config.Recursive {
		cmds = append(cmds, updateSubmodules(p.Config.SubmoduleRemote))
	}

	for _, cmd := range cmds {
		cmd.Dir = p.Build.Path
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

// Creates an empty git repository.
func initGit() *exec.Cmd {
	return exec.Command(
		"git",
		"init",
	)
}

// Sets the remote origin for the repository.
func remote(remote string) *exec.Cmd {
	return exec.Command(
		"git",
		"remote",
		"add",
		"origin",
		remote,
	)
}

// Checkout executes a git checkout command.
func checkoutHead() *exec.Cmd {
	return exec.Command(
		"git",
		"checkout",
		"-qf",
		"FETCH_HEAD",
	)
}

// Checkout executes a git checkout command.
func checkoutSha(commit string) *exec.Cmd {
	return exec.Command(
		"git",
		"reset",
		"--hard",
		"-q",
		commit,
	)
}

// fetch retuns git command that fetches from origin. If tags is true
// then tags will be fetched.
func fetch(ref string, tags bool, depth int) *exec.Cmd {
	tags_option := "--no-tags"
	if tags {
		tags_option = "--tags"
	}
	cmd := exec.Command(
		"git",
		"fetch",
		tags_option,
	)
	if depth != 0 {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--depth=%d", depth))
	}
	cmd.Args = append(cmd.Args, "origin")
	cmd.Args = append(cmd.Args, fmt.Sprintf("+%s:", ref))
	return cmd
}

// updateSubmodules recursively initializes and updates submodules.
func updateSubmodules(remote bool) *exec.Cmd {
	cmd := exec.Command(
		"git",
		"submodule",
		"update",
		"--init",
		"--recursive",
	)

	if remote {
		cmd.Args = append(cmd.Args, "--remote")
	}

	return cmd
}

// skipVerify returns a git command that, when executed configures git to skip
// ssl verification. This should may be used with self-signed certificates.
func skipVerify() *exec.Cmd {
	return exec.Command(
		"git",
		"config",
		"--global",
		"http.sslVerify",
		"false",
	)
}

// remapSubmodule returns a git command that, when executed configures git to
// remap submodule urls.
func remapSubmodule(name, url string) *exec.Cmd {
	name = fmt.Sprintf("submodule.%s.url", name)
	return exec.Command(
		"git",
		"config",
		name,
		url,
	)
}

// trace writes each command to stdout before it is executed. This is useful
// debugging the build to determine which commands were executed.
func trace(cmd *exec.Cmd) {
	fmt.Printf("<command>%s</command>\n", strings.Join(cmd.Args, " "))
}

// writeNetrc writes the netrc file to the user home directory.
func writeNetrc(machine, login, password string) error {
	if machine == "" {
		return nil
	}
	out := fmt.Sprintf(
		netrcFile,
		machine,
		login,
		password,
	)
	home := "/root"
	u, err := user.Current()
	if err == nil {
		home = u.HomeDir
	}
	path := filepath.Join(home, ".netrc")
	return ioutil.WriteFile(path, []byte(out), 0600)
}

// isDirEmpty returns true if the directory is empty. This is used to determine
// if the .git repository has been initialized yet.
func isDirEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return true
	}
	defer f.Close()

	_, err = f.Readdir(1)
	return err == io.EOF
}

// isPullRequest returns true if the event is a pull request event.
func isPullRequest(event string) bool {
	return event == "pull_request"
}

// isTag returns true if the event is a tag event.
func isTag(event, ref string) bool {
	return event == "tag" ||
		strings.HasPrefix(ref, "refs/tags/")
}

//
//
//
//
//

func main2() {

	app := cli.NewApp()
	app.Name = "git"
	app.Usage = "git clone plugin"
	app.Action = run
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
		// cli.Flag{
		// 	Name:   "submodule-override",
		// 	Desc:   "json map of submodule overrides",
		// 	EnvVar: "SUBMODULE_OVERRIDE",
		// },
	}
	app.Run(os.Args)

}

func run(c *cli.Context) {
	// plugin := Plugin{}
	// plugin.Repo.Clone = c.String("remote")
	// plugin.Build.Commit = c.String("sha")
	// plugin.Build.Event = c.String("event")
	// plugin.Build.Number = c.Int("number")
	// plugin.Build.Path = c.String("path")
	// plugin.Build.Ref = c.String("ref")
	// plugin.Netrc.Login = c.String("netrc.login")
	// plugin.Netrc.Machine = c.String("netrc.machine")
	// plugin.Netrc.Password = c.String("netrc.password")
	// plugin.Config.Depth = c.Int("depth")
	// plugin.Config.Recursive = c.BoolT("recursive")
	// plugin.Config.SkipVerify = c.Bool("skip-verify")
	// plugin.Config.SubmoduleRemote = c.Bool("submodule-override")
	// // plugin.Config.Submodules = c.String("submodule-override")

	// if err := plugin.Exec(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	_ = Plugin2{
		Repo: Repo{
			Clone: c.String("remote"),
		},
		Build: Build{
			Commit: c.String("sha"),
			Event:  c.String("event"),
			Number: c.Int("number"),
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
		},
	}

}

type (
	Repo struct {
		Clone string
	}

	Build struct {
		Path   string
		Event  string
		Number int
		Commit string
		Ref    string
	}

	Netrc struct {
		Machine  string
		Login    string
		Password string
	}

	Config struct {
		Depth           int
		Recursive       bool
		SkipVerify      bool
		Tags            bool
		Submodules      map[string]string
		SubmoduleRemote bool
	}

	Plugin2 struct {
		Repo   Repo
		Build  Build
		Netrc  Netrc
		Config Config
	}
)
