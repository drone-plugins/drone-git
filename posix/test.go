package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func testPullRequest(t *testing.T) {
	remote := "/tmp/remote/greeting"

	base, err := ioutil.TempDir("test", "")
	if err != nil {
		t.Error(err)
		return
	}

	for i, test := range tests {
		local := filepath.Join(base, fmt.Sprint(i))
		err = os.MkdirAll(local, 0644)
		if err != nil {
			t.Error(err)
			return
		}

		cmd := exec.Command("clone-commit")
		cmd.Env = []string{
			fmt.Sprintf("DRONE_BRANCH=%s", test.branch),
			fmt.Sprintf("DRONE_COMMIT=%s", test.commit),
			fmt.Sprintf("DRONE_WORKSPACE=%s", local),
			fmt.Sprintf("DRONE_REMOTE_URL=%s", remote),
		}

		_, err := cmd.CombinedOutput()
		if err != nil {
			t.Error(err)
			return
		}

		commit, err := getCommit(local)
		if err != nil {
			t.Error(err)
			return
		}

		branch, err := getCommit(local)
		if err != nil {
			t.Error(err)
			return
		}

		fmt.Println("remote")
		fmt.Println("    branch", test.branch)
		fmt.Println("    commit", test.commit)
		fmt.Println("  local:", local)
		fmt.Println("    branch", branch)
		fmt.Println("    commit", commit)

		// 1. test the commit sha
		// 2. test the branch name
		// 3. test the file content
	}
}

func getBranch(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func getCommit(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

var tests = []struct {
	branch string
	commit string
	tag    string
	file   string
	text   string
}{
	{
		commit: "51edb5943b8decc8ee57a1eef948576de8548fd3",
		branch: "master",
		tag:    "v1.0.0",
		file:   "hello.txt",
		text:   "hi world",
	},
	{
		commit: "da3a76a2f7e9174fcd1501c7022305f3ae7fa64b",
		branch: "master",
		tag:    "v1.1.0",
		file:   "hello.txt",
		text:   "hello world",
	},
	{
		commit: "277206e20ed92bdfe2c63b10d7bdf14ec8cd4c98",
		branch: "fr",
		tag:    "v2.0.0",
		file:   "hello.txt",
		text:   "salut monde",
	},
	{
		commit: "6692d23737d8cfa1f975fe7efc7e7dd3681898d7",
		branch: "fr",
		tag:    "v2.1.0",
		file:   "hello.txt",
		text:   "bonjour monde",
	},
}

//
// below we setup a local git directory with
// repeatable activity that can be used for testing.
//

// func setup() (name string, err error) {
// 	name, err = ioutil.TempDir("test", "")
// 	if err != nil {
// 		return
// 	}
// 	cmd := exec.Command("git", "init")
// 	cmd.Dir = name
// 	err = cmd.Run()
// 	if err != nil {
// 		return
// 	}

// 	for _, change := range changes {
// 		fn := filepath.Join(name, change.File)
// 		err = ioutil.WriteFile(fn, []byte(change.Data), 0644)
// 		if err != nil {
// 			return
// 		}
// 		cmd := exec.Command("git", "add", fn)
// 		cmd.Dir = name
// 		err = cmd.Run()
// 		if err != nil {
// 			return
// 		}

// 		cmd := exec.Command("git", "commit", "-m", "")
// 		cmd.Dir = name
// 		err = cmd.Run()
// 		if err != nil {
// 			return
// 		}
// 	}
// 	return
// }
