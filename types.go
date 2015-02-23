package main

import (
	"encoding/json"
	"os"
)

type Input struct {
	Clone struct {
		Origin string `json:"origin"`
		Remote string `json:"remote"`
		Branch string `json:"branch"`
		Sha    string `json:"sha"`
		Ref    string `json:"ref"`
		Dir    string `json:"dir"`

		Netrc struct {
			Machine  string `json:"machine"`
			Login    string `json:"login"`
			Password string `json:"user"`
		}

		Keypair struct {
			Public  string `json:"public"`
			Private string `json:"private"`
		}
	} `json:"clone"`

	User struct {
		Remote string `json:"remote"`
		Login  string `json:"login"`
		Name   string `json:"name"`
		Email  string `json:"email,omitempty"`
	} `json:"user"`

	Repo struct {
		Remote string `json:"remote"`
		Host   string `json:"host"`
		Owner  string `json:"owner"`
		Name   string `json:"name"`
		URL    string `json:"url"`
	} `json:"repo"`

	Commit struct {
		Status      string `json:"status"`
		Started     int64  `json:"started_at"`
		Finished    int64  `json:"finished_at"`
		Duration    int64  `json:"duration"`
		Sha         string `json:"sha"`
		Branch      string `json:"branch"`
		PullRequest string `json:"pull_request"`
		Author      string `json:"author"`
		Gravatar    string `json:"gravatar"`
		Timestamp   string `json:"timestamp"`
		Message     string `json:"message"`
	} `json:"commit"`

	Config struct {
		Image    string   `json:"image"`
		Env      []string `json:"env"`
		Script   []string `json:"script"`
		Branches []string `json:"branches"`
		Services []string `json:"services"`
	} `json:"config"`

	Params map[string]interface{} `json:"data"`
}

func Parse() (Input, error) {
	if len(os.Args) > 1 {
		return ParseArg()
	} else {
		return ParseStdin()
	}
}

func ParseMust() Input {
	in, err := Parse()
	if err != nil {
		panic(err)
	}
	return in
}

func ParseArg() (Input, error) {
	return ParseString(os.Args[1])
}

func ParseStdin() (Input, error) {
	in := Input{}
	err := json.NewDecoder(os.Stdin).Decode(&in)
	return in, err
}

func ParseString(raw string) (Input, error) {
	in := Input{}
	err := json.Unmarshal([]byte(raw), &in)
	return in, err
}
