
## Build

Build the project

```sh
make deps
make build
```

Creates Docker image `plugins/drone-git`

```sh
sudo make docker
```

## Usage

Clone a commit

```sh
./drone-git <<EOF
{
	"clone" : {
		"branch": "master",
		"remote": "git://github.com/drone/drone",
		"dir": "/drone/src/github.com/drone/drone",
		"ref": "refs/heads/master",
		"sha": "436b7a6e2abaddfd35740527353e78a227ddcb2c"
	}
}
EOF
```

Clone a pull request

```sh
./drone-git <<EOF
{
	"clone" : {
		"branch": "master",
		"remote": "git://github.com/drone/drone",
		"dir": "/drone/src/github.com/drone/drone",
		"ref": "refs/pull/892/merge",
		"sha": "8d6a233744a5dcacbf2605d4592a4bfe8b37320d"
	}
}
EOF
```

Clone a tag

```sh
./drone-git <<EOF
{
	"clone" : {
		"branch": "master",
		"remote": "git://github.com/drone/drone",
		"dir": "/drone/src/github.com/drone/drone",
		"sha": "339fb92b9629f63c0e88016fffb865e3e1055483",
		"ref": "refs/tags/v0.2.0"
	}
}
EOF
```

## Docker

Build the Docker container:

```sh
docker build -t plugins/drone-git .
```

Clone a repository inside the Docker container:

```sh
docker run -i plugins/drone-git <<EOF
{
	"clone" : {
		"branch": "master",
		"remote": "git://github.com/drone/drone",
		"dir": "/drone/src/github.com/drone/drone",
		"ref": "refs/heads/master",
		"sha": "436b7a6e2abaddfd35740527353e78a227ddcb2c"
	}
}
EOF
```
