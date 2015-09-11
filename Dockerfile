# Docker image for Drone's git-clone plugin
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-git .

FROM alpine:3.2
RUN apk-install ca-certificates git openssh curl perl
ADD drone-git /bin/
ENTRYPOINT ["/bin/drone-git"]
