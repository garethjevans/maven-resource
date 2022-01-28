FROM golang:1.17.5 as go
RUN GO111MODULES=on go get -u -ldflags="-s -w" github.com/garethjevans/maven-resource

FROM ubuntu:20.04
COPY --from=go /go/bin/maven-resource /bin/maven-resource

RUN maven-resource --help
# /opt/resource/in
# /opt/resource/out
