package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/drone/drone-plugin-go/plugin"
)

func main() {
	c := new(plugin.Clone)
	r := new(plugin.Repo)
	b := new(plugin.Build)
	plugin.Param("clone", c)
	plugin.Param("repo", r)
	plugin.Param("build", b)
	plugin.Parse()

	err := os.MkdirAll(c.Dir, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var cmds []*exec.Cmd
	if isPR(c) {
		cmds = append(cmds, clone(c))
		cmds = append(cmds, fetch(c))
		cmds = append(cmds, checkoutHead(c))
	} else {
		cmds = append(cmds, cloneBranch(c))
		cmds = append(cmds, checkoutSha(c))
	}

	for _, cmd := range cmds {
		cmd.Dir = c.Dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)
		err := cmd.Run()
		if err != nil {
			os.Exit(1)
		}
	}
}

// Returns true if cloning a pull request.
func isPR(b *plugin.Build) bool {
	return b.PullRequest.Number == 0
}

func isTag(b *plugin.Build) bool {
	return strings.HasPrefix(b.Commit.Ref, "refs/tags/")
}

// Clone executes a git clone command.
func clone(c *plugin.Clone, b *plugin.Build) *exec.Cmd {
	return exec.Command(
		"git",
		"clone",
		"--depth=50",
		"--recursive",
		b.Commit.Remote.Clone,
		c.Dir,
	)
}

// CloneBranch executes a git clone command
// for a single branch.
func cloneBranch(c *plugin.Clone, b *plugin.Build) *exec.Cmd {
	return exec.Command(
		"git",
		"clone",
		"-b",
		b.Commit.Ref,
		"--depth=50",
		"--recursive",
		b.Commit.Remote.Clone,
		c.Dir,
	)
}

// Checkout executes a git checkout command.
func checkoutSha(c *plugin.Clone) *exec.Cmd {
	return exec.Command(
		"git",
		"checkout",
		"-qf",
		c.Sha,
	)
}

// Checkout executes a git checkout command.
func checkoutHead(c *plugin.Clone) *exec.Cmd {
	return exec.Command(
		"git",
		"checkout",
		"-qf",
		"FETCH_HEAD",
	)
}

// Fetch executes a git fetch to origin.
func fetch(c *plugin.Clone) *exec.Cmd {
	return exec.Command(
		"git",
		"fetch",
		"origin",
		fmt.Sprintf("+%s:", c.Ref),
	)
}

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}
