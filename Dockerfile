FROM       golang:1.4
MAINTAINER paulcull <paul@pkhome.co.uk>

RUN        apt-get -qy update && apt-get -qy install vim-common gcc mercurial supervisor

WORKDIR    /go/src/github.com/paulcull/mqtt-webbrick
ADD        . /go/src/github.com/paulcull/mqtt-webbrick

RUN        go get -v

RUN  go build -ldflags " \
       -X main.buildVersion  $(grep "const Version " version.go | sed -E 's/.*"(.+)"$/\1/' ) \
       -X main.buildRevision $(git rev-parse --short HEAD) \
       -X main.buildBranch   $(git rev-parse --abbrev-ref HEAD) \
       -X main.buildDate     $(date +%Y%m%d-%H:%M:%S) \
       -X main.goVersion     $GOLANG_VERSION \
     "

EXPOSE     9980
CMD ./mqtt-webbrick
