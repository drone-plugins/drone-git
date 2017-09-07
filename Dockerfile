FROM alpine:3.5

RUN apk add --no-cache \
    ca-certificates \
    git \
    openssh \
    curl \
    perl

ADD drone-git /bin/
ENTRYPOINT ["/bin/drone-git"]
