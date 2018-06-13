FROM alpine:3.7
RUN apk add --no-cache ca-certificates git git-lfs openssh curl perl sudo

ADD posix/* /usr/local/bin/
RUN adduser -g Drone -s /bin/sh -D -u 1000 drone
RUN echo 'drone ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/drone
USER drone:drone
ENTRYPOINT ["/usr/local/bin/clone"]
