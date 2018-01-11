package main

import (
	"encoding/json"
	"time"
)

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
		Attempts        int
	}

	Backoff struct {
		Attempts int
		Duration time.Duration
	}
)

// below are special types used for unmarshaling structured data
// from environment variable or command line args.

type MapFlag struct {
	parts map[string]string
}

func (m *MapFlag) Get() map[string]string {
	return m.parts
}

func (m *MapFlag) Set(value string) error {
	m.parts = map[string]string{}
	return json.Unmarshal([]byte(value), &m.parts)
}

func (m *MapFlag) String() (s string) {
	return
}
