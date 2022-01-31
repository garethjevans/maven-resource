.PHONY: build
build:
	go build -o build/maven-resource main.go

test:
	go test ./...

lint:
	golangci-lint run

check-first-attempt: build
	cat cmd/testdata/check-first-attempt.json | ./build/maven-resource check

check: build
	cat cmd/testdata/check.json | ./build/maven-resource check

in: build
	mkdir -p test-output
	cat cmd/testdata/in.json | ./build/maven-resource in test-output
	tree test-output
	rm -fr test-output

test-all: test check check-first-attempt in
