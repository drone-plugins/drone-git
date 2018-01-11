package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Plugin struct {
	Repo    Repo
	Build   Build
	Netrc   Netrc
	Config  Config
	Backoff Backoff
}

func (p Plugin) Exec() error {
	var err error
	if p.Config.Attempts > 1 {
		fmt.Println("We will do up to", p.Config.Attempts,"attempts")
	}
	for i := 0; i < p.Config.Attempts; i++ {
		err = p.ExecActual()
		if err == nil {
			return nil
		}
		fmt.Println("Global retry attempt", i)
		os.RemoveAll(p.Build.Path)
	}
	return err
}

func (p Plugin) ExecActual() error {
	if p.Build.Path != "" {
		err := os.MkdirAll(p.Build.Path, 0777)
		if err != nil {
			return err
		}
	}

	err := writeNetrc(p.Netrc.Machine, p.Netrc.Login, p.Netrc.Password)
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
		buf := new(bytes.Buffer)
		cmd.Dir = p.Build.Path
		cmd.Stdout = io.MultiWriter(os.Stdout, buf)
		cmd.Stderr = io.MultiWriter(os.Stderr, buf)
		trace(cmd)
		err := cmd.Run()
		switch {
		case err != nil && shouldRetry(buf.String()):
			err = retryExec(cmd, p.Backoff.Duration, p.Backoff.Attempts)
			if err != nil {
				return err
			}
		case err != nil:
			return err
		}
	}

	return nil
}

// shouldRetry returns true if the command should be re-executed. Currently
// this only returns true if the remote ref does not exist.
func shouldRetry(s string) bool {
	return strings.Contains(s, "find remote ref")
}

// retryExec is a helper function that retries a command.
func retryExec(cmd *exec.Cmd, backoff time.Duration, retries int) (err error) {
	for i := 0; i < retries; i++ {
		// signal intent to retry
		fmt.Printf("retry in %v\n", backoff)

		// wait 5 seconds before retry
		<-time.After(backoff)

		// copy the original command
		retry := exec.Command(cmd.Args[0], cmd.Args[1:]...)
		retry.Dir = cmd.Dir
		retry.Stdout = os.Stdout
		retry.Stderr = os.Stderr
		trace(retry)
		err = retry.Run()
		if err == nil {
			return
		}
	}
	return
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
	tagsOption := "--no-tags"
	if tags {
		tagsOption = "--tags"
	}
	cmd := exec.Command(
		"git",
		"fetch",
		tagsOption,
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
		"--global",
		name,
		url,
	)
}
