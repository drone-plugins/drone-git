# Docker image for Drone's git-clone plugin
#
#     go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-git .

FROM gliderlabs/alpine:3.1
RUN apk-install ca-certificates git openssh curl perl
ADD drone-git /bin/
ENTRYPOINT ["/bin/drone-git"]