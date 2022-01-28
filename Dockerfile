FROM golang:1.17.5 as go
RUN GO111MODULES=on go install github.com/garethjevans/maven-resource

FROM scratch

# /opt/resource/in
# /opt/resource/out
