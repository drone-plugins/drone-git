FROM alpine:3.6

RUN apk update && \
  apk add \
    ca-certificates \
    git \
    openssh \
    curl \
    perl && \
  rm -rf /var/cache/apk/*

ADD drone-git /bin/
ENTRYPOINT ["/bin/drone-git"]
