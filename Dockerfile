# Docker image for Drone's git-clone plugin
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/git .

FROM alpine:3.2
RUN apk add -U ca-certificates git openssh curl perl && rm -rf /var/cache/apk/*

ENV LFS_VERSION 1.2.0
RUN curl -sLO https://github.com/github/git-lfs/releases/download/v${LFS_VERSION}/git-lfs-linux-amd64-${LFS_VERSION}.tar.gz && \
    tar xzf /git-lfs-linux-amd64-${LFS_VERSION}.tar.gz -C / && \
    mv /git-lfs-${LFS_VERSION}/git-lfs /usr/local/bin/ && \
    git-lfs init && \
    rm -rf /git-lfs-${LFS_VERSION} && \
    rm -rf /git-lfs-linux-amd64-${LFS_VERSION}.tar.gz

ADD drone-git /bin/
ENTRYPOINT ["/bin/drone-git"]
