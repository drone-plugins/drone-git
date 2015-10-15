package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/drone/drone-plugin-go/plugin"
)

// commits is a list of commits of different types (push, pull request, tag)
// to help us verify that this clone plugin can handle multiple commit types.
var commits = []struct {
	path   string
	clone  string
	event  string
	branch string
	commit string
	ref    string
	file   string
	data   string
}{
	// first commit
	{
		path:   "octocat/Hello-World",
		clone:  "https://github.com/octocat/Hello-World.git",
		event:  plugin.EventPush,
		branch: "master",
		commit: "553c2077f0edc3d5dc5d17262f6aa498e69d6f8e",
		ref:    "refs/heads/master",
		file:   "README",
		data:   "Hello World!",
	},
	// head commit
	{
		path:   "octocat/Hello-World",
		clone:  "https://github.com/octocat/Hello-World.git",
		event:  plugin.EventPush,
		branch: "master",
		commit: "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
		ref:    "refs/heads/master",
		file:   "README",
		data:   "Hello World!\n",
	},
	// pull request commit
	{
		path:   "octocat/Hello-World",
		clone:  "https://github.com/octocat/Hello-World.git",
		event:  plugin.EventPull,
		branch: "master",
		commit: "553c2077f0edc3d5dc5d17262f6aa498e69d6f8e",
		ref:    "refs/pull/208/merge",
		file:   "README",
		data:   "Goodbye World!\n",
	},
	// branch
	{
		path:   "octocat/Hello-World",
		clone:  "https://github.com/octocat/Hello-World.git",
		event:  plugin.EventPush,
		branch: "test",
		commit: "b3cbd5bbd7e81436d2eee04537ea2b4c0cad4cdf",
		ref:    "refs/heads/test",
		file:   "CONTRIBUTING.md",
		data:   "## Contributing\n",
	},
	// tags
	{
		path:   "github/mime-types",
		clone:  "https://github.com/github/mime-types.git",
		event:  plugin.EventTag,
		branch: "master",
		commit: "553c2077f0edc3d5dc5d17262f6aa498e69d6f8e",
		ref:    "refs/tags/v1.17",
		file:   ".gitignore",
		data:   "*.swp\n*~\n.rake_tasks~\nhtml\ndoc\npkg\npublish\ncoverage\n",
	},
}

// TestClone tests the ability to clone a specific commit into
// a fresh, empty directory every time.
func TestClone(t *testing.T) {

	for _, c := range commits {
		dir := setup()

		r := &plugin.Repo{Clone: c.clone}
		b := &plugin.Build{Commit: c.commit, Branch: c.branch, Ref: c.ref, Event: c.event}
		w := &plugin.Workspace{Path: dir}
		v := &Params{}
		if err := clone(r, b, w, v); err != nil {
			t.Errorf("Expected successful clone. Got error. %s.", err)
		}

		data := readFile(dir, c.file)
		if data != c.data {
			t.Errorf("Expected %s to contain [%s]. Got [%s].", c.file, c.data, data)
		}

		teardown(dir)
	}
}

// TestCloneNonEmpty tests the ability to clone a specific commit into
// a non-empty directory. This is useful if the git workspace is cached
// and re-stored for every build.
func TestCloneNonEmpty(t *testing.T) {
	dir := setup()
	defer teardown(dir)

	for _, c := range commits {

		r := &plugin.Repo{Clone: c.clone}
		b := &plugin.Build{Commit: c.commit, Branch: c.branch, Ref: c.ref, Event: c.event}
		w := &plugin.Workspace{Path: filepath.Join(dir, c.path)}
		v := &Params{}
		if err := clone(r, b, w, v); err != nil {
			t.Errorf("Expected successful clone. Got error. %s.", err)
			break
		}

		data := readFile(w.Path, c.file)
		if data != c.data {
			t.Errorf("Expected %s to contain [%s]. Got [%s].", c.file, c.data, data)
			break
		}
	}
}

// helper function that will setup a temporary workspace.
// to which we can clone the repositroy
func setup() string {
	dir, _ := ioutil.TempDir("/tmp", "drone_git_test_")
	os.Mkdir(dir, 0777)
	return dir
}

// helper function to delete the temporary workspace.
func teardown(dir string) {
	os.RemoveAll(dir)
}

// helper function to read a file in the temporary worskapce.
func readFile(dir, file string) string {
	filename := filepath.Join(dir, file)
	data, _ := ioutil.ReadFile(filename)
	return string(data)
}
