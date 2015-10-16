package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
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
	tags   []string
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
		tags:   nil,
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
		tags:   nil,
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
		tags:   nil,
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
		tags:   nil,
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
		tags: []string{
			"v1.16",
			"v1.17",
			"v1.17.1",
			"v1.17.2",
			"v1.18",
			"v1.19",
			"v1.20",
			"v1.20.1",
			"v1.21",
			"v1.22",
			"v1.23",
		},
	},
}

// TestClone tests the ability to clone a specific commit into
// a fresh, empty directory every time.
func TestClone(t *testing.T) {

	for _, c := range commits {
		dir := setup()

		tags := false
		if c.tags != nil {
			tags = true
		}

		r := &plugin.Repo{Clone: c.clone}
		b := &plugin.Build{Commit: c.commit, Branch: c.branch, Ref: c.ref, Event: c.event}
		w := &plugin.Workspace{Path: dir}
		v := &Params{Tags: tags}
		if err := clone(r, b, w, v); err != nil {
			t.Errorf("Expected successful clone. Got error. %s.", err)
		}

		data := readFile(dir, c.file)
		if data != c.data {
			t.Errorf("Expected %s to contain [%s]. Got [%s].", c.file, c.data, data)
		}

		if c.tags != nil {
			tags, err := getTags(dir)
			if err != nil {
				t.Error(err)
			}
			for _, tag := range c.tags {
				if !tags[tag] {
					t.Errorf("Expected tag [%s] to exist.", tag)
				}
			}
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

// TestClone tests if the arguments to `git fetch` are constructed properly.
func TestFetch(t *testing.T) {
	testdata := []struct {
		build *plugin.Build
		tags  bool
		depth int
		exp   []string
	}{
		{
			&plugin.Build{Ref: "refs/heads/master"},
			false,
			50,
			[]string{
				"git",
				"fetch",
				"--no-tags",
				"--depth=50",
				"origin",
				"+refs/heads/master:",
			},
		},
		{
			&plugin.Build{Ref: "refs/heads/master"},
			true,
			100,
			[]string{
				"git",
				"fetch",
				"--tags",
				"--depth=100",
				"origin",
				"+refs/heads/master:",
			},
		},
	}
	for _, td := range testdata {
		c := fetch(td.build, td.tags, td.depth)
		if len(c.Args) != len(td.exp) {
			t.Errorf("Expected: %s, got %s", td.exp, c.Args)
		}
		for i := range c.Args {
			if c.Args[i] != td.exp[i] {
				t.Errorf("Expected: %s, got %s", td.exp, c.Args)
			}
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

// getTags returns all of the tags in a git repository as a map.
func getTags(dir string) (map[string]bool, error) {
	cmd := exec.Command("git", "tag")
	cmd.Dir = dir
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	tags := make(map[string]bool)
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		tags[scanner.Text()] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return tags, nil
}
