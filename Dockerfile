# Docker image for Drone's git-clone plugin
#
#     docker build -t drone-plugins/drone-git .

FROM library/golang:1.4

# copy the local package files to the container's workspace.
ADD . /go/src/github.com/drone-plugins/drone-git/

# build the git-clone plugin inside the container.
RUN go get github.com/drone-plugins/drone-git/... && \
    go install github.com/drone-plugins/drone-git

# run the git-clone plugin when the container starts
ENTRYPOINT ["/go/bin/drone-git"]
