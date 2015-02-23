# docker build -t drone/drone-clone-git .

FROM library/golang:1.4.0
ADD . /gopath/src/github.com/drone/drone-clone-git/

WORKDIR /gopath/src/github.com/drone/drone-clone-git

RUN make && make install

ENTRYPOINT ["/gopath/src/github.com/drone/drone-clone-git/drone-clone-git"]
