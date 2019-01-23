FROM plugins/base:linux-amd64

LABEL maintainer="Drone.IO Community <drone-dev@googlegroups.com>" \
  org.label-schema.name="Drone Git" \
  org.label-schema.vendor="Drone.IO Community" \
  org.label-schema.schema-version="1.0"

RUN apk add --no-cache git openssh curl perl

ADD release/linux/amd64/drone-git /bin/
ENTRYPOINT ["/bin/drone-git"]
