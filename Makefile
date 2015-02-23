all: build

deps:
	go get

build:
	go build

docker:
	docker build --force-rm=false -t drone/drone-clone-git .

install:
	install -t /usr/local/bin drone-clone-git