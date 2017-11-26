FROM golang:1.6 as build

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app
RUN go get github.com/Sirupsen/logrus
RUN go get github.com/joho/godotenv
RUN go get github.com/urfave/cli
RUN env CGO_ENABLED=0 go build -ldflags "-s -w -X main.build=0" -v -a -tags netgo

FROM alpine:3.5

RUN apk update && \
  apk add \
  ca-certificates \
  git \
  openssh \
  curl \
  perl && \
  rm -rf /var/cache/apk/*

COPY --from=build /app/drone-git /bin/

ENTRYPOINT ["/bin/drone-git"]
