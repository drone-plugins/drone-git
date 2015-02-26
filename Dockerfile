# Docker image for Drone's git-clone plugin
#
#     docker build -t drone/drone-clone-git .

FROM library/golang:1.4

# copy the local package files to the container's workspace.
ADD . /go/src/github.com/drone/drone-clone-git/

# build the git-clone plugin inside the container.
RUN go get github.com/drone/drone-clone-git/... && \
    go install github.com/drone/drone-clone-git

# run the git-clone plugin when the container starts
ENTRYPOINT /go/bin/drone-clone-git
