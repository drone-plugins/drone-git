package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Plugin struct {
	Repo   Repo
	Build  Build
	Netrc  Netrc
	Config Config
}

func (p Plugin) Exec() error {
	if p.Build.Path != "" {
		err := os.MkdirAll(p.Build.Path, 0777)
		if err != nil {
			return err
		}
	}

	err := writeNetrc(p.Netrc.Login, p.Netrc.Login, p.Netrc.Password)
	if err != nil {
		return err
	}

	var cmds []*exec.Cmd

	if p.Config.SkipVerify {
		cmds = append(cmds, skipVerify())
	}

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
