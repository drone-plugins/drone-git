# Docker image for Drone's git-clone plugin
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/git .

FROM alpine:3.2
RUN apk add -U ca-certificates git openssh curl perl && rm -rf /var/cache/apk/*
ADD drone-git /bin/
ENTRYPOINT ["/bin/drone-git"]
