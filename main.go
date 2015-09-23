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

	"github.com/drone/drone-plugin-go/plugin"
)

var netrcFile = `
machine %s
login %s
password %s
`

func main() {
	v := struct {
		Depth int `json:"depth"`
	}{}

	r := new(plugin.Repo)
	b := new(plugin.Build)
	w := new(plugin.Workspace)
	plugin.Param("repo", r)
	plugin.Param("build", b)
	plugin.Param("workspace", w)
	plugin.Param("vargs", &v)
	plugin.MustParse()

	if v.Depth == 0 {
		v.Depth = 50
	}

	err := os.MkdirAll(w.Path, 0777)
	if err != nil {
		fmt.Printf("Error creating directory %s. %s\n", w.Path, err)
		os.Exit(2)
	}

	// generate the .netrc file
	if err := writeNetrc(w); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

	// write the rsa private key if provided
	if err := writeKey(w); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}

	var cmds []*exec.Cmd
	// check for a .git directory and whether it's empty
	if isDirEmpty(filepath.Join(w.Path, ".git")) {
		cmds = append(cmds, initGit())
		cmds = append(cmds, remote(r))
	}

	cmds = append(cmds, fetch(b, v.Depth))

	if isPR(b) {
		cmds = append(cmds, checkoutHead(b))
	} else {
		cmds = append(cmds, checkoutSha(b))
	}

	for _, cmd := range cmds {
		cmd.Dir = w.Path
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
	return strings.HasPrefix(b.Commit.Ref, "refs/pull/")
}

func isTag(b *plugin.Build) bool {
	return strings.HasPrefix(b.Commit.Ref, "refs/tags/")
}

// Creates an empty git repository.
func initGit() *exec.Cmd {
	return exec.Command(
		"git",
		"init",
	)
}

// Sets the remote origin for the repository.
func remote(r *plugin.Repo) *exec.Cmd {
	return exec.Command(
		"git",
		"remote",
		"add",
		"origin",
		r.Clone,
	)
}

// Checkout executes a git checkout command.
func checkoutSha(b *plugin.Build) *exec.Cmd {
	return exec.Command(
		"git",
		"checkout",
		"-qf",
		b.Commit.Sha,
	)
}

// Checkout executes a git checkout command.
func checkoutHead(b *plugin.Build) *exec.Cmd {
	return exec.Command(
		"git",
		"checkout",
		"-qf",
		"FETCH_HEAD",
	)
}

// Fetch executes a git fetch to origin.
func fetch(b *plugin.Build, depth int) *exec.Cmd {
	return exec.Command(
		"git",
		"fetch",
		fmt.Sprintf("--depth=%d", depth),
		"origin",
		fmt.Sprintf("+%s:", b.Commit.Ref),
	)
}

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}

// Writes the netrc file.
func writeNetrc(in *plugin.Workspace) error {
	if in.Netrc == nil || len(in.Netrc.Machine) == 0 {
		return nil
	}
	out := fmt.Sprintf(
		netrcFile,
		in.Netrc.Machine,
		in.Netrc.Login,
		in.Netrc.Password,
	)
	home := "/root"
	u, err := user.Current()
	if err == nil {
		home = u.HomeDir
	}
	path := filepath.Join(home, ".netrc")
	return ioutil.WriteFile(path, []byte(out), 0600)
}

// Writes the RSA private key
func writeKey(in *plugin.Workspace) error {
	if in.Keys == nil || len(in.Keys.Private) == 0 {
		return nil
	}
	home := "/root"
	u, err := user.Current()
	if err == nil {
		home = u.HomeDir
	}
	sshpath := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshpath, 0700); err != nil {
		return err
	}
	confpath := filepath.Join(sshpath, "config")
	privpath := filepath.Join(sshpath, "id_rsa")
	ioutil.WriteFile(confpath, []byte("StrictHostKeyChecking no\n"), 0700)
	return ioutil.WriteFile(privpath, []byte(in.Keys.Private), 0600)
}

func isDirEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return true
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false
}
