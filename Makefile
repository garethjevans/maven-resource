.PHONY: build
build:
	go build -o build/maven-resource main.go

test:
	go test ./...

lint:
	golangci-lint run

check: build
	cat check.json | ./build/maven-resource check

in: build
	cat in.json | ./build/maven-resource in
