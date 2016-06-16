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
)

// trace writes the command in the programs stdout for debug purposes.
// the command is wrapped in xml tags for easy parsing.
func trace(cmd *exec.Cmd) {
	fmt.Printf("+ %s\n", strings.Join(cmd.Args, " "))
}

// helper function returns true if directory dir is empty.
func isDirEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return true
	}
	defer f.Close()

	_, err = f.Readdir(1)
	return err == io.EOF
}

// helper function returns true if the commit is a pull_request.
func isPullRequest(event string) bool {
	return event == "pull_request"
}

// helper function returns true if the commit is a tag.
func isTag(event, ref string) bool {
	return event == "tag" ||
		strings.HasPrefix(ref, "refs/tags/")
}

// helper function to write a netrc file.
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

const netrcFile = `
machine %s
login %s
password %s
`
