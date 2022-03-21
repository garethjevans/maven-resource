FROM golang:1.18.0 as go
RUN GOPROXY=direct GO111MODULES=on go install -u -ldflags="-s -w" github.com/garethjevans/maven-resource@main

FROM ubuntu:20.04
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean

COPY --from=go /go/bin/maven-resource /bin/maven-resource

RUN maven-resource --help

COPY scripts/in /opt/resource/in
COPY scripts/check /opt/resource/check
