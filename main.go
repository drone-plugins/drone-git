package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	in := ParseMust()

	os.MkdirAll(in.Clone.Dir, 0700)

	var cmds []*exec.Cmd
	if isPR(in) || isTag(in) {
		cmds = append(cmds, clone(in))
		cmds = append(cmds, fetch(in))
		cmds = append(cmds, checkoutHead(in))
	} else {
		cmds = append(cmds, cloneBranch(in))
		cmds = append(cmds, checkoutSha(in))
	}

	for _, cmd := range cmds {
		cmd.Dir = in.Clone.Dir
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
func isPR(in Input) bool {
	return strings.HasPrefix(in.Clone.Ref, "refs/pull/")
}

func isTag(in Input) bool {
	return strings.HasPrefix(in.Clone.Ref, "refs/tags/")
}

// Clone executes a git clone command.
func clone(in Input) *exec.Cmd {
	return exec.Command(
		"git",
		"clone",
		"--depth=50",
		"--recursive",
		in.Clone.Remote,
		in.Clone.Dir,
	)
}

// CloneBranch executes a git clone command
// for a single branch.
func cloneBranch(in Input) *exec.Cmd {
	//branch := fmt.Sprintf("--branch=%s", in.Clone.Branch)
	return exec.Command(
		"git",
		"clone",
		"-b",
		in.Clone.Branch,
		"--depth=50",
		"--recursive",
		in.Clone.Remote,
		in.Clone.Dir,
	)
}

// Checkout executes a git checkout command.
func checkoutSha(in Input) *exec.Cmd {
	return exec.Command(
		"git",
		"checkout",
		"-qf",
		in.Clone.Sha,
	)
}

// Checkout executes a git checkout command.
func checkoutHead(in Input) *exec.Cmd {
	return exec.Command(
		"git",
		"checkout",
		"-qf",
		"FETCH_HEAD",
	)
}

// Fetch executes a git fetch to origin.
func fetch(in Input) *exec.Cmd {
	return exec.Command(
		"git",
		"fetch",
		"origin",
		fmt.Sprintf("+%s:", in.Clone.Ref),
	)
}

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}
